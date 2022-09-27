package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"PackageServer/cache"
	"PackageServer/constant"
	"PackageServer/logger"
	"PackageServer/util"
)

var (
	ErrTokenNotFound   error = errors.New("token not found")
	ErrTokenExpired    error = errors.New("Token expired")
	ErrTokenAccessFail error = errors.New("Token access fail")
	ErrTokenInValid    error = errors.New("Token not valid")
	ErrTokenMalformed  error = errors.New("That's not even a token")

	SignKey string
)

// actually we implement session mode.
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			logger.Log.Infof("middleware:JwtAuth: %v", ErrTokenNotFound.Error())
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    ErrTokenNotFound.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}
		logger.Log.Info("middleware:JwtAuth: token in header: ", token)

		j := util.NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			logger.Log.Infof("middleware:JwtAuth: %+v", err)
			if err == ErrTokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": -1,
					"msg":    ErrTokenExpired.Error(),
					"data":   nil,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    errors.Cause(err).Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}
		logger.Log.Infof("middleware:JwtAuth: claims is %v", claims)

		jwtCache := cache.NewJsonWebTokenCache()
		err = jwtCache.ValidateJwt(claims.User, token, string(constant.From))
		if err != nil {
			logger.Log.Infof("middleware:JwtAuth: %v", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    errors.Cause(err).Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}
		logger.Log.Info("middleware:JwtAuth: jwt validate pass")
		c.Next()

	}
}
