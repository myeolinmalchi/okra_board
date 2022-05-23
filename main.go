package main

import (
	"io"
	"log"
	"okra_board/config"
	"okra_board/controllers"
	mw "okra_board/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){

    logFile, err := config.InitLoggingFile()
    if err != nil {
        log.Println("로그 파일을 생성하지 못했습니다.")
        log.Println(err.Error())
        return
    }

    logger := log.New(logFile, "INFO: ", log.LstdFlags)

    if err := config.InitDBConnection(); err != nil {
        logger.Println("데이터베이스에 접근할 수 없습니다.")
        logger.Println(err.Error())
        return
    }

    gin.DefaultWriter = io.MultiWriter(logFile)

    config, err := config.LoadConfig()
    if err != nil {
        logger.Println("설정 파일을 찾지 못했습니다. 서버를 종료합니다.")
        return
    } 

    os.Setenv("SECRET_KEY", config.SecretKey)

    whiteList := make(map[string]bool)
    for _, ipaddr := range config.WhiteList {
        whiteList[ipaddr] = true
    }
    IPWhiteList := mw.IPWhiteList(whiteList)

    gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    r.Use(cors.New(cors.Config {
        AllowAllOrigins:    true,
        AllowMethods:       []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:       []string{"Content-Type", "Authorization"},
        ExposeHeaders:      []string{"Authorization"},
        AllowCredentials:   true,
        MaxAge: 12 * time.Hour,
    }))
    r.Static("/images", "./public/images")
    v1 := r.Group("/api/v1")
    {
        /* route for Posts */
        v1.GET("posts", controllers.GetPosts)
        v1.GET("posts/:postId", controllers.GetPost)

        v1.POST("posts", IPWhiteList, mw.Auth(), controllers.InsertPost)
        v1.PUT("posts/:postId", IPWhiteList, mw.Auth(), controllers.UpdatePost)
        v1.DELETE("posts/:postId", IPWhiteList, mw.Auth(), controllers.DeletePost)
        v1.POST("posts/select",IPWhiteList , mw.Auth(), controllers.ResetSelectedPost)

        v1.GET("thumbnails", controllers.GetThumbnails)
        v1.GET("thumbnails/selected", controllers.GetSelectedThumbnails)
        /* route for Posts */

        /* route for Admin */
        v1.POST("admin", mw.Auth(), controllers.AdminUpsert(true))
        v1.PUT("admin", mw.Auth(), controllers.AdminUpsert(false))
        v1.POST("admin/login", controllers.AdminLogin)
        /* route for Admin */
    }
    r.Run(":3000")
}
