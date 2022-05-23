package middleware

import (
	"fmt"
	"okra_board/services"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Headers", "Content-Type, Autorization, Origin")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT, OPTIONS")
        c.Header("Access-Control-Request-Methods", "GET, DELETE, POST, PUT, OPTIONS")
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

/*
func Auth2() gin.HandlerFunc {
    return func(c *gin.Context) {
        accessToken := c.Request.Header.Get("AccessToken")
        if accessToken == "" {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "access token is empty",
            })
            c.Abort()
        } else {
            err, isExpired := services.VerifyAccessToken(accessToken)
            if err != nil {
                c.JSON(401, gin.H {
                    "status": 401,
                    "message": "invalid access token",
                })
                c.Abort()
            } else if isExpired {
                verifyRefreshToken(c)
            }
        }
    }
}

func verifyRefreshToken(c *gin.Context) {
    refreshToken := c.Request.Header.Get("RefreshToken")
    if refreshToken == "" {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "refresh token is empty",
        })
        c.Abort()
    } else {
        err, isExpired := services.VerifyRefreshToken(refreshToken)
        if err != nil {
            c.JSON(401, gin.H {
                "status": 401,
                "message": "invalid refresh token",
            })
            c.Abort()
        } else if isExpired {
            
        }
    }
}

func verifyRefreshToken2(c *gin.Context) {
    refreshToken := c.Request.Header.Get("RefreshToken")
    if refreshToken == "" {
        c.JSON(401, gin.H {
            "status": 401,
            "message": "refresh token is empty",
        })
        c.Abort()
    } else {
        err, claims := services.GetTokenClaims(refreshToken)
        if err != nil {

        } else {
            id := claims["id"]
        }
    }
}
*/

func IPWhiteList(whiteList map[string]bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !whiteList[c.ClientIP()] {
            fmt.Println("Permission denied: "+c.ClientIP())
            c.AbortWithStatusJSON(403, gin.H {
                "status": 403,
                "message": "Permission denied",
            })
        }
    }
}

