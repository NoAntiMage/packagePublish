package util

import "os"

func IsExist(src string) bool {
	_, err := os.Stat(src)
	return err == nil || os.IsExist(err)
}
