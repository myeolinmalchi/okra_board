package main

import (
    "./config"
    "./services"
    "./controllers"
    "github.com/gin-gonic/gin"
)

func main(){

    config.InitDBConnection()
    gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    r.Use(CORSMiddleware())

    v1 := r.Group("/api/v1")
    {
        v1.GET("posts", controllers.GetPosts)
        v1.GET("posts/:postId", controllers.GetPost)
        v1.POST("posts", AuthMiddleware(), controllers.InsertPost)
        v1.PUT("posts/:postId", AuthMiddleware(), controllers.UpdatePost)
        v1.DELETE("posts/:postId", AuthMiddleware(), controllers.DeletePost)
        v1.POST("posts/select", AuthMiddleware(), controllers.ResetSelectedPost)

        v1.GET("thumbnails", controllers.GetThumbnails)
        v1.GET("thumbnails/selected", controllers.GetSelectedThumbnails)

        v1.POST("admin", AuthMiddleware(), controllers.AdminRegist)
        v1.POST("admin/login", controllers.AdminLogin)
    }
    r.Run(":3000")
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Headers", "Content-Type, Autorization, Origin")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT")
    }
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.Request.Header.Get("Authorization")
        if token == "" {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "token is empty",
            })
            c.Abort()
        } else if err := services.VerifyToken(token); err != nil {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "invalid token",
            })
            c.Abort()
        }
    }
}

