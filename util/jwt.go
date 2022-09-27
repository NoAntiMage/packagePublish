package util

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"PackageServer/constant"
	"PackageServer/logger"
)

var (
	ErrTokenAccessFail error = errors.New("Token access fail")
	ErrTokenIsInvalid  error = errors.New("Token is invalid")
	SignKey            string
)

type CustomClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

type Jwt struct {
	SigningKey []byte
}

func NewJwt() Jwt {
	return Jwt{
		SigningKey: []byte(InitSignKey()),
	}
}

func InitSignKey() string {
	SignKey = constant.Secret
	return SignKey
}

func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *Jwt) CreateTokenWithExpire(claims CustomClaims) (string, error) {
	expireAt := time.Now().Add(time.Hour * 1).Unix()
	claims.StandardClaims.ExpiresAt = expireAt
	return j.CreateToken(claims)
}

func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "util:Jwt:")
	}

	logger.Log.Debugf("parsed Jwt : %v", token)

	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, errors.Wrap(ErrTokenIsInvalid, "util:Jwt:")
}

func (j *Jwt) UpdateToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil

	})
	if err != nil {
		return "", errors.Wrap(err, "util:Jwt:")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		return j.CreateTokenWithExpire(*claims)
	}
	return "", errors.Wrap(ErrTokenAccessFail, "util:Jwt:")
}
