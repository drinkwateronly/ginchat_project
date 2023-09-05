package utils

import (
	"errors"
	"ginchat/define"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaim struct {
	Account              string
	Username             string
	jwt.RegisteredClaims // 不要写成RegisteredClaims jwt.RegisteredClaims
}

func GenerateJWT(account, name string) (string, error) {
	//
	uc := UserClaim{
		Account:  account,
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "J. C.",
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(define.TokenExpireTime))),
		},
	}
	// t包含三部分：Header Claims Signature
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// 加密token，并生成str
	signedStr, err := t.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return signedStr, nil
}

func ValidateJwt(token string) (*UserClaim, error) {
	// 新建userClaim结构体
	uc := new(UserClaim)
	// jwt.ParseWithClaims 输入 需要解析的JWT字符串、一个实现了jwt.Claims接口的结构体、用于提供验证签名所需的密钥JwtKey的回调函数
	claims, err := jwt.ParseWithClaims(
		token,
		uc,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(define.JwtKey), nil
		})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, errors.New("token is invalid")
	}
	return uc, err
}
