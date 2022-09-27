package service

import (
	"PackageServer/config"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/repo"
	"PackageServer/util"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var ErrNotProperArea = errors.New("ErrNotProperArea")
var ErrPackageNameFormat = errors.New("ErrPackageNameFormat")
var ErrFileChunkNotMatch = errors.New("ErrFileChunkNotMatch")
var ErrEnvChunkNotMatch = errors.New("ErrEnvChunkNotMatch")
var ErrWrongEnvFile = errors.New("ErrWrongEnvFile")
var ErrLackOfChunks = errors.New("ErrLackOfChunks")
var ErrPackMd5NotMatch = errors.New("ErrPackMd5NotMatch")
var ErrPackageNotFound = errors.New("ErrpackageNotFound")

type PackageManagementSvc interface {
	GetPackageInfo(do dto.PublishPlanLog) (*dto.PackageInfo, error)
	PackageExist(packageInfo dto.PackageInfo) error
	GetPackageMd5(packageInfo dto.PackageInfo) (string, error)
	PackInfoReceive(postInfo dto.PackageInfoPost) error
	ChunkUpload(c *gin.Context, fileHeader *multipart.FileHeader, chunkInfo dto.ChunkInfo) error
	PackCheck(packageName string) (*dto.ChunkLack, error)
	PackTmpGarbageCollect(path string, fileName string, chunkNum int)
}

func NewPackageManagementSvc() PackageManagementSvc {
	return &packageMgtSvc{
		areaInfoRepo:      repo.NewAreaInfoRepo(),
		serviceOnlineRepo: repo.NewServiceOnlineRepo(),
	}
}

type packageMgtSvc struct {
	areaInfoRepo      repo.AreaInfoRepo
	serviceOnlineRepo repo.ServiceOnlineRepo
}

func (pm *packageMgtSvc) GetPackageInfo(do dto.PublishPlanLog) (*dto.PackageInfo, error) {
	var packInfo dto.PackageInfo
	areaMo, err := pm.areaInfoRepo.GetById(do.AreaInfoId)
	if err != nil {
		return nil, err
	}

	serviceMo, err := pm.serviceOnlineRepo.GetById(do.ServiceOnlineId)
	if err != nil {
		return nil, err
	}
	nowDate := time.Now().Format("20060102")
	packInfo.FileName = fmt.Sprintf("%v_tag_%v_%v.%v", serviceMo.ServiceName, areaMo.AreaName, do.Version, serviceMo.ArchiveType)
	packInfo.Path = fmt.Sprintf("%vtag_%v_%v/", config.ServerConf.PackageDir, areaMo.AreaName, nowDate)

	if err := pm.PackageExist(packInfo); err != nil {
		return nil, err
	}
	md5Secret, err := pm.GetPackageMd5(packInfo)
	if err != nil {
		return nil, err
	}
	packInfo.Md5Sum = md5Secret
	return &packInfo, nil
}

func (pm *packageMgtSvc) PackageExist(packageInfo dto.PackageInfo) error {
	fileLocation := fmt.Sprintf("%v%v", packageInfo.Path, packageInfo.FileName)
	if _, err := os.Stat(fileLocation); err != nil {
		logger.Log.Errorf("%v", err)
		return errors.Wrap(ErrPackageNotFound, "service:packageMgtSvc:")
	}
	return nil
}

func (pm *packageMgtSvc) GetPackageMd5(packageInfo dto.PackageInfo) (string, error) {
	fileLocation := fmt.Sprintf("%v%v", packageInfo.Path, packageInfo.FileName)
	m := util.NewMd5Util()
	return m.FileMd5Sum(fileLocation)
}

func (pm *packageMgtSvc) PackInfoReceive(postInfo dto.PackageInfoPost) error {
	receiver, err := pm.parsePostInfo(postInfo)
	if err != nil {
		return err
	}

	if err = pm.IsProperArea(*receiver); err != nil {
		return err
	}

	if err = pm.makeUpdateDir(*receiver); err != nil {
		return err
	}

	if err = pm.savePackInfoToEnv(*receiver); err != nil {
		return err
	}
	return nil
}

func (pm *packageMgtSvc) ChunkUpload(c *gin.Context, chunk *multipart.FileHeader, chunkInfo dto.ChunkInfo) error {
	logger.Log.Debugf("PackageManagementSvc:ChunkUpload: chunk info : %v", chunkInfo)

	if pm.isChunkOfFile(chunkInfo.ChunkName, chunkInfo.FileName) == false {
		return errors.Wrapf(ErrFileChunkNotMatch, "PackageManagementSvc:ChunkUpload")
	}
	packInfoEnvDto, err := pm.loadPackInfoFromEnv(chunkInfo.FileName)
	if err != nil {
		return err
	}

	if err = pm.chunkVerifyWithPackInfo(chunkInfo, packInfoEnvDto); err != nil {
		return err
	}

	if err = pm.saveChunk(c, chunk, chunkInfo.ChunkName); err != nil {
		return err
	}

	return nil
}

/* description:
as a worker, after package info received and chunks uploaded,
check if all chunks have been received.
merge all chunks to package and md5 it.
*/
func (pm *packageMgtSvc) PackCheck(packageName string) (*dto.ChunkLack, error) {
	var do dto.ChunkLack
	do.FileName = packageName
	packInfoEnv, err := pm.loadPackInfoFromEnv(packageName)
	if err != nil {
		return nil, err
	}

	lackList, err := pm.chunksCount(packInfoEnv)
	if err != nil {
		return nil, err
	}
	do.LackList = lackList
	if len(lackList) != 0 {
		return &do, errors.Wrapf(ErrLackOfChunks, "packageMgtSvc:PackCheck")
	}

	if err := pm.mergeChunkToFile(packInfoEnv); err != nil {
		return nil, err
	}

	if err := pm.PackMd5Check(packInfoEnv); err != nil {
		return nil, err
	}

	pm.PackTmpGarbageCollect(pm.getUpdateDir(), packInfoEnv.FileName, packInfoEnv.ChunkNum)

	return nil, nil
}

func (pm *packageMgtSvc) chunksCount(do dto.PackageInfoEnv) (lackList []int, err error) {
	updateDir := pm.getUpdateDir()
	files, err := ioutil.ReadDir(updateDir)
	if err != nil {
		return nil, errors.Wrapf(err, "packageMgtSvc:chunksCount")
	}
	var chunkIdList []int
	for _, f := range files {
		name := f.Name()
		if pm.isChunkOfFile(name, do.FileName) {
			chunkId, err := pm.getChunkId(name)
			if err != nil {
				return nil, err
			}
			chunkIdList = append(chunkIdList, chunkId)
		}
	}
	sort.Ints(chunkIdList)
	logger.Log.Debugf("packageMgtSvc:chunksCount: %v exist chunkId: %v", do.FileName, chunkIdList)

	// double pointer for 2 arrays, lackList and ChunkNum(List)
	var i = 0
	var p2 = 1
	for p2 < do.ChunkNum+1 {
		if i >= len(chunkIdList) {
			break
		}
		p1 := chunkIdList[i]
		if p1 != p2 {
			lackList = append(lackList, p2)
		} else {
			i++
		}
		p2++
	}
	// add rest from ChunkNum(List) to LackLast
	for p2 < do.ChunkNum+1 {
		lackList = append(lackList, p2)
		p2++
	}

	logger.Log.Debugf("packageMgtSvc:chunksCount: %v lack chunkId: %v", do.FileName, lackList)
	return lackList, nil
}

func (pm *packageMgtSvc) mergeChunkToFile(do dto.PackageInfoEnv) error {
	updateDir := pm.getUpdateDir()

	merger := util.FileMerger{
		WorkPath: updateDir,
		DistFile: do.FileName,
		DistPath: updateDir,
	}
	return merger.ChunkToFile(do.ChunkNum)
}

func (pm *packageMgtSvc) PackMd5Check(do dto.PackageInfoEnv) error {
	packLocation := pm.getUpdateDir() + do.FileName
	md5Tool := util.NewMd5Util()
	newMd5Secret, err := md5Tool.FileMd5Sum(packLocation)
	if err != nil {
		return err
	}
	logger.Log.Debugf("packageMgtSvc:PackMd5Check: NewMd5: %v", newMd5Secret)
	if newMd5Secret != do.Md5Sum {
		logger.Log.Infof("packageMgtSvc:PackMd5Check: %v: info: %v,new: %v", ErrPackMd5NotMatch.Error(), do.Md5Sum, newMd5Secret)
		return errors.Wrapf(ErrPackMd5NotMatch, "packageMgtSvc:PackMd5Check")
	}
	return nil
}

func (pm *packageMgtSvc) PackTmpGarbageCollect(path string, fileName string, chunkNum int) {
	for i := 1; i < chunkNum+1; i++ {
		ChunkLoc := util.GetChunkLocation(path, fileName, i)
		if err := os.Remove(ChunkLoc); err == nil {
			logger.Log.Debugf("packageMgtSvc:PackTmpGarbageCollect: gc %v", ChunkLoc)
		}
	}
}

func (pm *packageMgtSvc) isChunkOfFile(chunkName string, fileName string) bool {
	l1 := strings.Split(chunkName, ".")
	if l1[len(l1)-1] != "tmp" {
		return false
	}
	l2 := strings.Split(l1[len(l1)-2], "_")
	if len(l2) != 2 && l2[0] != "chunk" {
		return false
	}
	l3 := l1[0 : len(l1)-2]
	ShouldBeFileName := strings.Join(l3, ".")
	if ShouldBeFileName != fileName {
		return false
	}
	return true
}

func (pm *packageMgtSvc) getChunkId(chunkName string) (int, error) {
	l1 := strings.Split(chunkName, ".")
	l2 := strings.Split(l1[len(l1)-2], "_")
	chunkId, err := strconv.Atoi(l2[1])
	return chunkId, errors.Wrapf(err, "packageMgtSvc:getChunkId")
}

func (pm *packageMgtSvc) loadPackInfoFromEnv(fileName string) (dto.PackageInfoEnv, error) {
	var do dto.PackageInfoEnv
	updateDir := pm.getUpdateDir()
	envFileLocation := pm.getEnvFileLocation(updateDir, fileName)
	envFile, err := os.OpenFile(envFileLocation, os.O_RDONLY, 0)
	defer envFile.Close()
	if err != nil {
		return do, errors.Wrapf(err, "PackageManagementSvc:loadPackInfoFromEnv")
	}
	//	logger.Log.Debugf("PackageManagementSvc: read from .env: %v", envFileLocation)

	m := make(map[string]interface{})
	br := bufio.NewReader(envFile)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		//	logger.Log.Debugf("PackageManagementSvc: read env line: %v", string(line))
		l := strings.Split(string(line), "=")
		if len(l) != 2 {
			return do, errors.Wrapf(ErrWrongEnvFile, "PackageManagementSvc:loadPackInfoFromEnv")
		}
		if l[0] == "ChunkNum" {
			m[l[0]], _ = strconv.Atoi(l[1])
			continue
		}
		m[l[0]] = l[1]
	}
	//	logger.Log.Debugf("PackageManagementSvc: m : %v", m)

	if err = util.Map2Struct(&do, m); err != nil {
		return do, errors.Wrapf(err, "PackageManagementSvc:loadPackInfoFromEnv")
	}

	validator := validator.New()
	if err := validator.Struct(do); err != nil {
		logger.Log.Debugf("PackageManagementSvc: Err PackageInfoEnv : %v", do)
		return do, errors.Wrapf(ErrWrongEnvFile, "PackageManagementSvc:loadPackInfoFromEnv")
	}
	logger.Log.Debugf("PackageManagementSvc: load env file data: %v", do)
	return do, nil
}

func (pm *packageMgtSvc) chunkVerifyWithPackInfo(chunkInfo dto.ChunkInfo, packInfoEnv dto.PackageInfoEnv) error {
	if chunkInfo.FileName != packInfoEnv.FileName {
		return errors.Wrapf(ErrEnvChunkNotMatch, "PackageManagementSvc:chunkVerifyWithPackInfo")
	}
	l := strings.Split(chunkInfo.ChunkName, ".")
	l2 := strings.Split(l[len(l)-2], "_")
	chunkId, err := strconv.Atoi(l2[1])
	if err != nil {
		return errors.Wrapf(err, "PackageManagementSvc:chunkVerifyWithPackInfo")
	}
	if chunkId > packInfoEnv.ChunkNum {
		return errors.Wrapf(ErrEnvChunkNotMatch, "PackageManagementSvc:chunkVerifyWithPackInfo")
	}
	return nil
}

func (pm *packageMgtSvc) saveChunk(c *gin.Context, chunk *multipart.FileHeader, chunkName string) error {
	updateDir := pm.getUpdateDir()
	chunkLocation := fmt.Sprintf("%v%v", updateDir, chunkName)
	logger.Log.Debugf("PackageManagementSvc:saveChunk: prepare to save %v", chunkLocation)
	return c.SaveUploadedFile(chunk, chunkLocation)
}

// example package name: spd-wms_tag_spd-local_20220824-0932.jar
func (pm *packageMgtSvc) parsePostInfo(postInfo dto.PackageInfoPost) (*dto.PackageInfoEnv, error) {
	allList := strings.Split(postInfo.FileName, ".")
	infoList := strings.Split(allList[0], "_")
	if len(infoList) != 4 || infoList[1] != "tag" {
		return nil, errors.Wrapf(ErrPackageNameFormat, "PackageManagementSvc: %v", postInfo.FileName)
	}
	var receiver = dto.PackageInfoEnv{
		FileName:    postInfo.FileName,
		ServiceName: infoList[0],
		AreaName:    infoList[2],
		Version:     infoList[3],
		Md5Sum:      postInfo.Md5Sum,
		ChunkNum:    postInfo.ChunkNum,
	}

	logger.Log.Debugf("PackageManagementSvc: receiver: %v", receiver)
	return &receiver, nil
}

func (pm *packageMgtSvc) IsProperArea(InfoReceiver dto.PackageInfoEnv) error {
	if InfoReceiver.AreaName != config.ServerConf.AreaName {
		errMsg := fmt.Sprintf("PackageManagementSvc: area is %v", config.ServerConf.AreaName)
		return errors.Wrapf(ErrNotProperArea, errMsg)
	}
	return nil
}

// mkdir ${workdir}/${date}
func (pm *packageMgtSvc) makeUpdateDir(InfoReceiver dto.PackageInfoEnv) error {
	updateDir := pm.getUpdateDir()
	if util.IsExist(updateDir) {
		return nil
	} else {
		if err := os.Mkdir(updateDir, 0775); err != nil {
			return errors.Wrapf(err, "PackageManagementSvc:")
		}
	}
	return nil
}

func (pm *packageMgtSvc) getEnvFileLocation(updateDir string, fileName string) string {
	return fmt.Sprintf("%v.%v-info", updateDir, fileName)
}

// touch ${workdir}/${date}/.{env}
func (pm *packageMgtSvc) savePackInfoToEnv(InfoReceiver dto.PackageInfoEnv) error {
	updateDir := pm.getUpdateDir()
	envFileLocation := pm.getEnvFileLocation(updateDir, InfoReceiver.FileName)
	envFile, err := os.OpenFile(envFileLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer envFile.Close()
	if err != nil {
		return errors.Wrapf(err, "PackageManagementSvc:")
	}
	logger.Log.Debugf("PackageManagementSvc: touch .env: %v", envFileLocation)

	m, err := util.Struct2Map(InfoReceiver)
	if err != nil {
		return err
	}
	for k, v := range m {
		line := fmt.Sprintf("%v=%v\n", k, v)
		_, err = envFile.WriteString(line)
		if err != nil {
			return errors.Wrapf(err, "PackageManagementSvc:")
		}
	}
	return nil
}

func (pm *packageMgtSvc) getUpdateDir() string {
	nowDate := time.Now().Format("0102")
	updateDir := fmt.Sprintf("%v%v/", config.ServerConf.PackageDir, nowDate)
	return updateDir
}
