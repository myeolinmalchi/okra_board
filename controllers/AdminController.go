package controllers

import (
    "okra_board/services"
    "okra_board/models"
    "github.com/gin-gonic/gin"
)

func AdminUpsert(isRegist bool) func(c *gin.Context) {
    return func(c *gin.Context) {
        requestBody := &models.Admin{}    
        err := c.ShouldBind(requestBody)
        if err != nil {
            c.JSON(400, err.Error())
            return
        }

        isValid, validationResult := services.AdminUpsert(isRegist, requestBody)
        if isValid {
            c.Status(200)
        } else if validationResult == nil {
            c.Status(400)
        } else {
            c.IndentedJSON(422, validationResult)
        }
    }
}

func AdminLogin(c *gin.Context) {
    requestBody := &models.Admin{}
    err := c.ShouldBind(requestBody)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    
    if services.Login(requestBody) {
        token, _ := services.CreateToken(requestBody.ID)
        c.Header("Authorization", token)
        c.Status(200)
    } else {
        c.Status(401)
    }
}
