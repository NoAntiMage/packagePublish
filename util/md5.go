package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/pkg/errors"
)

type Md5Util interface {
	Md5Sum(string) string
	FileMd5Sum(string) (string, error)
}

type md5Util struct{}

func NewMd5Util() Md5Util {
	return &md5Util{}
}

func (md *md5Util) Md5Sum(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func (md *md5Util) FileMd5Sum(fileLocation string) (string, error) {
	pFile, err := os.Open(fileLocation)
	if err != nil {
		return "", errors.Wrapf(err, "md5Util:FileMd5Sum")
	}
	defer pFile.Close()
	m := md5.New()
	io.Copy(m, pFile)
	return hex.EncodeToString(m.Sum(nil)), nil
}
