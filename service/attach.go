package service

import (
	"fmt"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func UpLoad(c *gin.Context) {
	writer := c.Writer
	req := c.Request
	file, header, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(writer, "")
	}
	suffix := ".png"
	fileName := header.Filename
	split := strings.Split(fileName, ".")
	if len(split) != 0 {
		suffix = "." + split[len(split)-1]
	}
	//byteF := make([]byte, header.Size)
	//file.Read(byteF)
	newFileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	// 存放路径
	url := "./asset/upload/" + newFileName
	saveFile, err := os.Create(url)
	if err != nil {
		utils.RespFail(writer, "")
	}
	_, err = io.Copy(saveFile, file)
	defer saveFile.Close()
	if err != nil {
		utils.RespFail(writer, "")
	}
	utils.RespOK(writer, url, "图片发送成功成功")
}
