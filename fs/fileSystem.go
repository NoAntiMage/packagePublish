package fs

import (
	"PackageServer/config"
	"PackageServer/logger"
	"errors"
	"os"
)

var (
	ErrCreatePackageDirFail = errors.New("ErrCreatePackageDirFail")
)

func InitFs() error {
	packageDir := config.ServerConf.PackageDir
	if packageDir == "" {
		logger.Log.Errorf("config init empty, %v", ErrCreatePackageDirFail)
		return ErrCreatePackageDirFail
	}

	if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
		logger.Log.Errorf("%v", ErrCreatePackageDirFail)
	}

	return nil
}
