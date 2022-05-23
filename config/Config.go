package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDBConnection() (err error){
    dsn := "root:382274@tcp(localhost:3306)/board_prototype?parseTime=true"
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            ti, _ := time.LoadLocation("Asia/Seoul")
            return time.Now().In(ti)
        },
    })
    return
}

type Config struct {
    WhiteList   []string    `json:"whitelist"`
    SecretKey   string      `json:"secretkey"`
}

func LoadConfig() (*Config, error){
    file, err := os.Open("config.json")
    defer file.Close()
    config := &Config{}
    jsonParser := json.NewDecoder(file)
    jsonParser.Decode(config)
    return config, err
}

func InitLoggingFile() (*os.File, error) {
    startTime := time.Now().Format("2006-01-02")
    logFile := "log/log-" + startTime
    return os.Create(strings.TrimSpace(logFile))
}
