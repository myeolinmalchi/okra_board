package main

import (
    "okra_board/config"
    "okra_board/controllers"
    "github.com/gin-gonic/gin"
    mw "okra_board/middleware"
)

func main(){

    config.InitDBConnection()
    gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    r.Use(mw.CORS())
    r.Static("/images", "./public/images")

    whiteList := make(map[string]bool)
    whiteList["124.54.16.127"] = true

    r.Use(mw.IPWhiteList(whiteList))

    v1 := r.Group("/api/v1")
    {
        /* route for Posts */
        v1.GET("posts", controllers.GetPosts)
        v1.GET("posts/:postId", controllers.GetPost)

        v1.POST("posts", mw.Auth(), controllers.InsertPost)
        v1.PUT("posts/:postId", mw.Auth(), controllers.UpdatePost)
        v1.DELETE("posts/:postId", mw.Auth(), controllers.DeletePost)

        v1.POST("posts/select", mw.Auth(), controllers.ResetSelectedPost)

        v1.GET("thumbnails", controllers.GetThumbnails)
        v1.GET("thumbnails/selected", controllers.GetSelectedThumbnails)
        /* route for Posts */

        /* route for Admin */
        v1.POST("admin", mw.Auth(), controllers.AdminRegist)
        v1.POST("admin/login", controllers.AdminLogin)
        /* route for Admin */
    }
    r.Run(":3000")
}
