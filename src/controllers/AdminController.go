package controllers

import (
    "../services"
    "../models"
    "github.com/gin-gonic/gin"
)

func AdminRegist(c *gin.Context) {
    requestBody := &models.Admin{}    
    err := c.ShouldBind(requestBody)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }

    isValid, validationResult := services.Regist(requestBody)
    if isValid {
        c.Status(200)
    } else if validationResult == nil {
        c.Status(400)
    } else {
        c.IndentedJSON(422, validationResult)
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
