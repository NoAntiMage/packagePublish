package service

import (
	"PackageServer/config"
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/repo"
	"PackageServer/variable"
	"time"
)

type PublishPlanSvc interface {
	Add(dto.PublishPlanLog) (id uint, err error)
	GetById(id uint) (*dto.PublishPlanLog, error)
	PublishVersion(do dto.PublishPlanLog) error
}

func NewPublishPlanSvc() PublishPlanSvc {
	return &publishPlanSvc{
		planStageSvc:    NewPlanStageSvc(),
		publishPlanRepo: repo.NewPublishPlanRepo(),
		areaInfoSvc:     NewAreaInfoSvc(),
		packageMgtSvc:   NewPackageManagementSvc(),
		rpcLoginSvc:     NewRpcLoginService(),
	}
}

type publishPlanSvc struct {
	planStageSvc    PlanStageSvc
	publishPlanRepo repo.PublishPlanRepo
	areaInfoSvc     AreaInfoSvc
	packageMgtSvc   PackageManagementSvc
	rpcLoginSvc     RpcLoginSvc
}

func (p *publishPlanSvc) Add(do dto.PublishPlanLog) (id uint, err error) {
	if do.FromDate.IsZero() {
		do.FromDate = time.Now()
	}
	return p.publishPlanRepo.Add(do)
}
func (p *publishPlanSvc) GetById(id uint) (*dto.PublishPlanLog, error) {
	mo, err := p.publishPlanRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	PublishPlanLogDto := &dto.PublishPlanLog{
		AreaInfoId:      mo.AreaInfoId,
		ServiceOnlineId: mo.ServiceOnlineId,
		Version:         mo.Version,
		FromDate:        mo.FromDate,
	}
	return PublishPlanLogDto, nil
}

/*
PlanStage
---
1. PlanCreated
|
2. PlanPublished
|
3. PackageExist - ErrPackageNotFound
|
4. RemoteAlive - ErrRemoteNotAlive
| mode: release
5. RpcLoginCheck -ErrRpcLoginFail
|
FileServerJob
*/

// 1. get package info
// 2. check remote server alive
// 3. start filePushJob and update PlanStage
func (p *publishPlanSvc) PublishVersion(planDto dto.PublishPlanLog) error {
	var packInfoDto *dto.PackageInfo
	var updation *dto.PlanStageUpdate
	var job *dto.PackagePushJob

	// 1. PlanCreated
	id, err := p.Add(planDto)
	if err != nil {
		return err
	}
	logger.Log.Debugf("service:publishPlan: add plan %v", planDto)

	// 2. PlanPublished
	updation = dto.PublishPlanLog2Updation(planDto)
	updation.Stage = constant.PlanPublished
	if err := p.planStageSvc.Update(updation); err != nil {
		return err
	}

	// 3. PackageExist
	packInfoDto, err = p.packageMgtSvc.GetPackageInfo(planDto)
	if err != nil {
		updation.Stage = constant.ErrPackageInvalid.Error()
		if err := p.planStageSvc.Update(updation); err != nil {
			return err
		}
		return err
	}
	logger.Log.Debugf("service:publishPlan: packageInfo %v ", *packInfoDto)

	// 4. RemoteAlive
	areaDto, err := p.areaInfoSvc.GetById(planDto.AreaInfoId)
	if err != nil {
		return err
	}

	err = p.areaInfoSvc.RemoteHealthCheck(areaDto.AreaName)
	if err != nil {
		updation.Stage = constant.ErrRemoteNotReady.Error()
		if err := p.planStageSvc.Update(updation); err != nil {
			return err
		}
		return err
	}

	logger.Log.Debugf("service:publishPlan: remoteUrl alive: %v", areaDto.AreaName)

	// 5. rpcTokenCheck - mode: release
	if config.ServerConf.RouterMode == "release" {
		if err = p.rpcLoginSvc.RemoteAreaTokenCheck(areaDto.AreaName); err != nil {
			updation.Stage = constant.ErrLoginTokenCheckFail.Error()
			if err := p.planStageSvc.Update(updation); err != nil {
				return err
			}
			return err
		}
	}

	updation.Stage = constant.PlanReadyToPush
	if err := p.planStageSvc.Update(updation); err != nil {
		return err
	}
	logger.Log.Debugf("service:publishPlan: PlanStage: %v", updation)

	remoteUrl, err := p.areaInfoSvc.GetUrl(areaDto.AreaName)
	if err != nil {
		return err
	}
	logger.Log.Debugf("service:publishPlan: remoteFtp: %v", remoteUrl)

	// FileServerJob
	job = &dto.PackagePushJob{
		Id:        id,
		FileName:  packInfoDto.FileName,
		Path:      packInfoDto.Path,
		Md5Sum:    packInfoDto.Md5Sum,
		Area:      areaDto.AreaName,
		RemoteUrl: remoteUrl,
	}
	variable.FilePushJobChan <- *job

	updation.Stage = constant.PlanPushing
	if err := p.planStageSvc.Update(updation); err != nil {
		return err
	}
	logger.Log.Debugf("service:publishPlan: PlanStage: %v", updation)

	return nil
}
