package config

import (
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
