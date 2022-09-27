package service

import (
	"PackageServer/config"
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/rpc"
	"PackageServer/util"

	"github.com/pkg/errors"
)

type FileTransferSvc interface {
	StartPackagePushJob(do dto.PackagePushJob) error
	UploadChunksWithRetry(do dto.PackagePushJob, chunkNum int) (failList []int, err error)
}

func NewFileTransferSvc() FileTransferSvc {
	return &fileTransferSvc{
		publishPlanSvc: NewPublishPlanSvc(),
		planStageSvc:   NewPlanStageSvc(),
		packMgtSvc:     NewPackageManagementSvc(),
	}
}

type fileTransferSvc struct {
	splitter       *util.FileSplitter
	publishPlanSvc PublishPlanSvc
	planStageSvc   PlanStageSvc
	packMgtSvc     PackageManagementSvc
}

/*
1. split package to chunks
2. post info to remote server
3. push chunks
4. validate remote server package

PlanStage
---
Pushing - ErrChunkUploadFail
|
Pushed
*/
func (f *fileTransferSvc) StartPackagePushJob(do dto.PackagePushJob) error {
	mo, err := f.publishPlanSvc.GetById(do.Id)
	if err != nil {
		return err
	}
	var updation = &dto.PlanStageUpdate{
		AreaInfoId:      mo.AreaInfoId,
		ServiceOnlineId: mo.ServiceOnlineId,
		Version:         mo.Version,
	}

	f.splitter = &util.FileSplitter{
		TargetFile:  do.FileName,
		SrcPath:     do.Path,
		WorkPath:    config.ServerConf.PackageDir,
		ChunkPrefix: do.FileName,
		ChunkSize:   1 * (1 << 20),
	}
	chunkNum, err := f.splitPackage()
	if err != nil {
		return err
	}

	packageInfoPost := &dto.PackageInfoPost{
		FileName: do.FileName,
		Md5Sum:   do.Md5Sum,
		ChunkNum: chunkNum,
	}

	if err := f.postPackageInfo(do.Area, do.RemoteUrl, packageInfoPost); err != nil {
		return err
	}
	logger.Log.Info("service:ftp:StartPackagePushJob: postPackInfo successfully.")

	failList, err := f.UploadChunksWithRetry(do, chunkNum)
	if err != nil {
		logger.Log.Infof("fileTransferSvc:StartPackagaPushJob: uploadFailList: %v", failList)
		return err
	}

	if err := f.packageCheck(do.Area, do.RemoteUrl, packageInfoPost.FileName); err != nil {
		updation.Stage = constant.ErrChunkUploadFail.Error()
		if err := f.planStageSvc.Update(updation); err != nil {
			logger.Log.Debugf("fileTransferSvc: PlanStage: %v", updation)
		}
		logger.Log.Infof("fileTransferSvc:StartPackagaPushJob: %+v", err)
		return err
	}
	logger.Log.Info("service:ftp:StartPackagePushJob: uploadChunks successfully.")

	updation.Stage = constant.PlanPushed
	if err := f.planStageSvc.Update(updation); err != nil {
		return err
	}
	logger.Log.Debugf("fileTransferSvc: PlanStage: %v", updation)

	f.packMgtSvc.PackTmpGarbageCollect(do.Path, do.FileName, chunkNum)
	logger.Log.Info("service:ftp:StartPackagePushJob: local tmp chunks have been sweeped.")

	return nil
}

func (f *fileTransferSvc) splitPackage() (num int, err error) {
	num, err = f.splitter.FileToChunk()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (f *fileTransferSvc) postPackageInfo(area string, remoteUrl string, info *dto.PackageInfoPost) error {
	packPublishRpc, err := rpc.NewPackagePublishRpc(area)
	if err != nil {
		return err
	}
	return packPublishRpc.PostPackageInfo(remoteUrl, info)
}

func (f *fileTransferSvc) UploadChunksWithRetry(do dto.PackagePushJob, chunkNum int) (failList []int, err error) {
	var chunkIdList []int
	for i := 1; i < chunkNum+1; i++ {
		chunkIdList = append(chunkIdList, i)
	}
	failList, _ = f.uploadChunks(do, chunkIdList)

	retryLimit := 3
	for len(failList) != 0 {
		if retryLimit == 0 {
			break
		}
		failList, _ = f.uploadChunks(do, failList)
		retryLimit--
	}

	if len(failList) != 0 {
		return failList, errors.Wrapf(constant.ErrChunkUploadFail, "fileTransferSvc:uploadChunksWithRetry")
	}
	return nil, nil
}

func (f *fileTransferSvc) uploadChunks(do dto.PackagePushJob, chunkIdList []int) (failList []int, err error) {
	for _, id := range chunkIdList {
		chunkLoc := util.GetChunkLocation(do.Path, do.FileName, id)
		if err := f.uploadChunk(do.Area, do.RemoteUrl, chunkLoc); err != nil {
			failList = append(failList, id)
			logger.Log.Debugf("fileTransferSvc:uploadChunks: %v", err)
			continue
		}
	}
	return
}

func (f *fileTransferSvc) uploadChunk(area string, remoteUrl string, chunkLocation string) error {
	logger.Log.Debugf("fileTransferSvc:uploadChunk: uploading: %v", chunkLocation)
	packPublishRpc, err := rpc.NewPackagePublishRpc(area)
	if err != nil {
		return err
	}
	return packPublishRpc.PostChunkUpload(remoteUrl, chunkLocation)
}

func (f *fileTransferSvc) packageCheck(area string, remoteUrl string, PackName string) error {
	packPublishRpc, err := rpc.NewPackagePublishRpc(area)
	if err != nil {
		return err
	}
	return packPublishRpc.GetPackCheck(remoteUrl, PackName)
}
