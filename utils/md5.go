package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmpStr := h.Sum([]byte(nil))
	return hex.EncodeToString(tmpStr)
}

//func MD5Encode(data string) string {
//	return strings.ToUpper(Md5Encode(data))
//}

func Md5EncodeByte(data []byte) string {
	h := md5.New()
	h.Write(data)
	tmpStr := h.Sum([]byte(nil))
	return hex.EncodeToString(tmpStr)
}

func MakeGroupId() string {
	groupId := ""
	groupIdLength := rand.Intn(3) + 7
	for i := 0; i < groupIdLength; i++ {
		tmpNum := strconv.Itoa(rand.Intn(9))
		groupId += tmpNum
	}
	return groupId
}

func MakePassword(rawPassword, salt string) string {
	return Md5Encode(rawPassword + salt)
}

func ValidatePassword(rawPassword, salt, Md5Password string) bool {
	return Md5Encode(rawPassword+salt) == Md5Password
}
