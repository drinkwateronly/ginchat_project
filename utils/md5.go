package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmpStr := h.Sum([]byte(nil))
	return hex.EncodeToString(tmpStr)
}

func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func MakePassword(rawPassword, salt string) string {
	return Md5Encode(rawPassword + salt)
}

func ValidatePassword(rawPassword, salt, Md5Password string) bool {
	return Md5Encode(rawPassword+salt) == Md5Password
}
