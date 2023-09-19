package utils

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var DB *gorm.DB
var RDB *redis.Client
var CTX *gin.Context

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

func InitRedis() {
	db, _ := strconv.Atoi(viper.GetString("redis.db"))
	poolSize, _ := strconv.Atoi(viper.GetString("redis.poolSize"))
	maxIdleConn, _ := strconv.Atoi(viper.GetString("redis.maxIdleConn"))

	RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           db,
		PoolSize:     poolSize,
		MaxIdleConns: maxIdleConn,
	})
}

const (
	PublishKey = "publish_key"
)

func Publish(ctx context.Context, channel string, message string) error {
	var err error
	err = RDB.Publish(ctx, channel, message).Err()
	return err
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RDB.Subscribe(ctx, channel)
	message, err := sub.ReceiveMessage(ctx)
	return message.Payload, err
}
