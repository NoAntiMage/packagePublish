package util

import (
	"PackageServer/logger"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	defaultChunkSize = 1 * (1 << 20)
)

var (
	ErrTargetFileNotFound error = errors.New("target file is not found")
)

type FileSplitter struct {
	TargetFile  string
	SrcPath     string
	WorkPath    string
	ChunkPrefix string
	ChunkSize   int
}

type FileMerger struct {
	WorkPath    string
	ChunkPrefix string
	DistFile    string
	DistPath    string
}

func (fs *FileSplitter) FileToChunk() (num int, err error) {
	fileLocation := fs.SrcPath + fs.TargetFile

	f, err := os.OpenFile(fileLocation, os.O_RDONLY, os.ModePerm)
	defer f.Close()
	if err != nil {
		return 0, errors.Wrap(err, "util:FileSplitter:fileToChunk")
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return 0, errors.Wrap(err, "util:FileSplitter:fileToChunk")
	}

	var fileSize int64 = fileInfo.Size()

	if fs.ChunkSize == 0 {
		fs.ChunkSize = defaultChunkSize
	}
	totalPartsNum := int(math.Ceil(float64(fileSize) / float64(fs.ChunkSize)))

	logger.Log.Infof("Split %v to %d pieces.\n", fs.TargetFile, totalPartsNum)

	for i := int(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(fs.ChunkSize), float64(fileSize-int64(i*fs.ChunkSize))))
		partBuffer := make([]byte, partSize)

		f.Read(partBuffer)

		chunkLocation := GetChunkLocation(fs.SrcPath, fs.TargetFile, i+1)

		chunk, err := os.Create(chunkLocation)

		if err != nil {
			return 0, errors.Wrap(err, "util:FileSplitter:fileToChunk")
		}

		chunk.Write(partBuffer)
		chunk.Close()

		logger.Log.Debugf("Split to %v \n", chunkLocation)
	}
	return totalPartsNum, nil
}

/*
description:
merge chunk to file.
if file has existed,backup it.
*/
func (fm *FileMerger) ChunkToFile(num int) error {
	distLocation := fm.DistPath + fm.DistFile
	_, err := os.Stat(distLocation)
	if err == nil {
		logger.Log.Infof("FileMerger:ChunkToFile: %v has existed. It will be backuped.", fm.DistFile)
		os.Rename(distLocation, fmt.Sprintf("%v.bak", distLocation))
	}
	fout, err := os.OpenFile(distLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer fout.Close()
	if err != nil {
		return errors.Wrapf(err, "util:FileSplitter:fileToChunk")
	}

	for i := int(0); i < num; i++ {
		chunkLocation := fmt.Sprintf("%v%v.chunk_%v.tmp", fm.WorkPath, fm.DistFile, strconv.Itoa((i + 1)))
		logger.Log.Debugf("util:FileMerger:chunkToFile: merge chunk %v", chunkLocation)
		f, err := os.OpenFile(chunkLocation, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "util:FileMerger:chunkToFile")
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.Wrapf(err, "util:FileMerger:chunkToFile")
		}
		fout.Write(b)
		f.Close()
	}
	logger.Log.Debugf("merged file : %v", distLocation)
	return nil
}

func GetChunkLocation(path string, fileName string, id int) string {
	return fmt.Sprintf("%v%v.chunk_%v.tmp", path, fileName, strconv.Itoa(id))
}
