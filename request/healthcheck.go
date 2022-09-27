package request

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// curl /ping
func GetPing(do dto.AreaInfo) error {
	url := fmt.Sprintf("http://%v:%v%vping", do.IpAddr, do.Port, do.UrlPath)
	logger.Log.Debugf("url %v", url)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "request:GetPing:")
	}

	rbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "request:GetPing:")
	}
	var baseRes dto.BaseResponse
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	logger.Log.Debugf("result: %v", baseRes)

	if baseRes.Msg == "pong" && baseRes.Code == 1 {
		return nil

	} else {
		return errors.Wrap(constant.AreaNotAlive, "request:GetPing:")
	}
}
