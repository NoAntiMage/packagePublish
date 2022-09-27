package request

import (
	"PackageServer/dto"
	"PackageServer/logger"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

// func PostPackInfo(targetUrl string, info *dto.PackageInfoPost) error {
func PostPackInfo(reqInfoDto dto.RequestInfo, info *dto.PackageInfoPost) error {
	var baseRes dto.BaseResponse
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", reqInfoDto.TargetUrl, payload)
	if err != nil {
		return errors.Wrapf(err, "request:PostPackInfo")
	}

	req.Header.Add("Content-Type", "application/json")
	if reqInfoDto.Token != "" {
		req.Header.Add("token", reqInfoDto.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request:PostPackInfo")
	}
	defer resp.Body.Close()
	rbyte, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return errors.Wrap(err, "json")
	}
	logger.Log.Debugf("Url: %v, Request: %v, Response: %v", reqInfoDto.TargetUrl, string(data), string(rbyte))
	return nil
}

func PostChunkUpload(reqInfoDto dto.RequestInfo, chunkLocation string) error {
	var baseRes dto.BaseResponse

	chunkInfoDto := dto.ChunkLocToChunkInfo(chunkLocation)

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	fileWriter, err := bodyWriter.CreateFormFile("chunk", chunkInfoDto.ChunkName)
	if err != nil {
		return errors.Wrapf(err, "request:PostChunkUpload")
	}

	f, err := os.Open(chunkLocation)
	defer f.Close()
	if err != nil {
		return errors.Wrapf(err, "request:PostChunkUpload")
	}
	io.Copy(fileWriter, f)

	bodyWriter.Close()

	req, err := http.NewRequest("POST", reqInfoDto.TargetUrl, bodyBuf)
	if err != nil {
		return errors.Wrapf(err, "request:PostChunkUpload")
	}
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Add("packageName", chunkInfoDto.FileName)

	if reqInfoDto.Token != "" {
		req.Header.Add("token", reqInfoDto.Token)
	}
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request:PostChunkUpload")
	}
	defer resp.Body.Close()

	rbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "request:PostChunkUpload")
	}

	logger.Log.Debugf("Url: %v, Response: %v", reqInfoDto.TargetUrl, string(rbyte))
	if err := json.Unmarshal(rbyte, &baseRes); err != nil {
		return errors.Wrap(err, "json")
	}
	if baseRes.Code != 1 {
		return errors.New(baseRes.Msg)
	}
	return nil
}

func GetPackCheck(reqInfoDto dto.RequestInfo, packName string) error {
	var baseRes dto.BaseResponse

	req, err := http.NewRequest("GET", reqInfoDto.TargetUrl, nil)
	if err != nil {
		return errors.Wrapf(err, "request:GetPackCheck")
	}
	query := req.URL.Query()
	query.Add("packageName", packName)
	req.URL.RawQuery = query.Encode()
	if reqInfoDto.Token != "" {
		req.Header.Add("token", reqInfoDto.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request:GetPackCheck")
	}
	defer resp.Body.Close()

	rbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "request:GetPackCheck")
	}
	logger.Log.Debugf("Url: %v, Response: %v .", reqInfoDto.TargetUrl, string(rbyte))
	err = json.Unmarshal(rbyte, &baseRes)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	if baseRes.Code != 1 {
		return errors.New(baseRes.Msg)
	}
	return nil
}
