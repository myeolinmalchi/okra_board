package middleware

import (
    "okra_board/services"
    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Headers", "Content-Type, Autorization, Origin")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT")
    }
}

func Auth() gin.HandlerFunc {
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

func IPWhiteList(whiteList map[string]bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !whiteList[c.ClientIP()] {
            c.AbortWithStatusJSON(403, gin.H {
                "status": 403,
                "message": "Permission denied",
            })
        }
    }
}

