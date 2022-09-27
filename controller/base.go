package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

type DefaultReponse struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

type SuccessReponse struct {
	DefaultReponse
	Data interface{} `json:"Data"`
}

func (b *BaseApi) Success(c *gin.Context, msg string, obj interface{}) {
	var res SuccessReponse
	res.Code = 1
	res.Msg = msg
	res.Data = obj

	c.JSON(http.StatusOK, res)
}

func (b *BaseApi) Error400(c *gin.Context, err error) {
	var res DefaultReponse
	res.Code = 0
	res.Msg = err.Error()

	c.JSON(http.StatusBadRequest, res)
}

func (b *BaseApi) Error404(c *gin.Context, err error) {
	var res DefaultReponse
	res.Code = 0
	res.Msg = err.Error()

	c.JSON(http.StatusNotFound, res)
}

func (b *BaseApi) Error502(c *gin.Context, err error) {
	var res DefaultReponse
	res.Code = 0
	res.Msg = err.Error()

	c.JSON(http.StatusBadGateway, res)
}
