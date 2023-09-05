package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config MySQL:", viper.Get("mysql"))
}

func InitMySQL() {
	// 自定义日志模板，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)
	var err error = nil
	// 若要处理err，目前认为只能提前先声明err，因为DB是全局变量，DB,err:=时，db就变成了局部变量
	DB, err = gorm.Open(mysql.Open(
		viper.GetString("mysql.dns")),
		&gorm.Config{
			Logger: newLogger, // log
		})
	if err != nil {
		panic("failed to connect database")
	}
}
