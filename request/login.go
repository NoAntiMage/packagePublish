package request

import (
	"PackageServer/config"
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var ErrConfigInvalid error = errors.New("config is invalid")
var ErrTokenUpdateFail error = errors.New("ErrTokenUpdateFail")

// @api: /timestamp?user=${user}
// @return timestamp
func GetTimeStamp(loginTokenDto dto.LoginToken, areaInfoDto dto.AreaInfo) (*dto.TimeStamp, error) {
	var timeStampDto dto.TimeStamp
	ip := areaInfoDto.IpAddr
	port := areaInfoDto.Port
	addr := fmt.Sprintf("%v:%v", ip, port)
	url := fmt.Sprintf("http://%v/timestamp?user=%v", addr, loginTokenDto.User)
	logger.Log.Debugf("area %v, url: %v", loginTokenDto.Area, url)

	if ip == "" || port == "" {
		return nil, errors.Wrap(ErrConfigInvalid, "request:login")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "request:GetTimeStamp")
	}

	rbyte, _ := ioutil.ReadAll(resp.Body)
	var baseRes dto.BaseResponse
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return nil, errors.Wrap(err, "json")
	}

	logger.Log.Debugf("result: %v", baseRes)

	if baseRes.Msg == "ok" && baseRes.Code == 1 {

		jsonStr, err := json.Marshal(baseRes.Data)
		if err != nil {
			return nil, errors.Wrap(err, "json")
		}
		err = json.Unmarshal([]byte(jsonStr), &timeStampDto)
		if err != nil {
			return nil, errors.Wrap(err, "json")
		}

		return &timeStampDto, nil
	} else {
		err = errors.New("502")
		return nil, err
	}
}

// @api: /digestToken
// @return jwt
func PostDigestToken(loginTokenDto dto.LoginToken, areaInfoDto dto.AreaInfo) (jwt string, err error) {
	ip := areaInfoDto.IpAddr
	port := areaInfoDto.Port
	addr := fmt.Sprintf("%v:%v", ip, port)
	url := fmt.Sprintf("http://%v%v/digestToken", addr, constant.UrlVersion)
	logger.Log.Debugf("area %v, url: %v", loginTokenDto.Area, url)

	if ip == "" || port == "" {
		return "", errors.Wrap(ErrConfigInvalid, "request:PostDigestToken:")
	}

	buf, err := json.Marshal(loginTokenDto)
	if err != nil {
		return "", errors.Wrap(err, "request:PostDigestToken:")
	}
	dataReader := bytes.NewReader(buf)

	resp, err := http.Post(url, "application/json", dataReader)
	if err != nil {
		return "", errors.Wrap(err, "request:PostDigestToken:")
	}
	defer resp.Body.Close()

	rbyte, _ := ioutil.ReadAll(resp.Body)
	var baseRes dto.BaseResponse
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return "", errors.Wrap(err, "json")
	}

	logger.Log.Debugf("result: %v", baseRes)

	if baseRes.Code == 0 {
		return "", errors.Wrap(errors.New(baseRes.Msg), "request:PostDigestToken:")
	}

	value := baseRes.Data["jwt"]
	v, ok := value.(string)
	if !ok {
		return "", errors.Wrapf(constant.ErrTypeAssertFail, "request:login")
	}
	jwt = v
	return jwt, nil
}

// @api /rpcToken/refreshExpire
// @return rpcToken
func PostRpcTokenExpireTime(reqInfoDto dto.RequestInfo, ttl int) error {
	logger.Log.Debugf("url: %v", reqInfoDto.TargetUrl)

	var rpcTokenUpdateDto = dto.RpcTokenUpdate{
		User:       config.ServerConf.Name,
		ExpireTime: ttl,
	}

	buf, err := json.Marshal(rpcTokenUpdateDto)
	if err != nil {
		return errors.Wrap(err, "request:PostRpcTokenExpireTime")
	}
	dataReader := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", reqInfoDto.TargetUrl, dataReader)
	if err != nil {
		return errors.Wrap(err, "request:PostRpcTokenExpireTime")
	}
	req.Header.Add("Content-Type", "application/json")
	if reqInfoDto.Token != "" {
		req.Header.Add("token", reqInfoDto.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request:PostRpcTokenExpireTime")
	}
	defer resp.Body.Close()

	rbyte, _ := ioutil.ReadAll(resp.Body)
	var baseRes dto.BaseResponse
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	logger.Log.Debugf("result: %v", baseRes)

	if baseRes.Msg == "ok" && baseRes.Code == 1 {
		return nil
	} else {
		return errors.Wrapf(ErrTokenUpdateFail, "request:PostRpcTokenExpireTime")
	}
}

// @api /rpcToken/del?area=${area}
func GetRpcTokenDelete() {}
