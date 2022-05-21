package controllers

import (
	"strconv"
	"okra_board/models"
	"okra_board/services"
    "okra_board/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosts(c *gin.Context) {
    size, err1 := strconv.Atoi(c.DefaultQuery("size", "15")) 
    page, err2 := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err1 != nil || err2 != nil {
        c.Status(400)
        return
    }
    query := config.DB.Table("posts")
    boardIdStr, boardIdExists := c.GetQuery("boardId")
    keyword, keywordExists := c.GetQuery("keyword")
    token := c.Request.Header.Get("Authorization")
    if err := services.VerifyToken(token); token == "" || err != nil {
        query = query.Where("status = ?", true)
    } 
    if boardIdExists {
        boardId, err := strconv.Atoi(boardIdStr)
        if err != nil {
            c.Status(400)
            return
        }
        query = query.Where("board_id = ?", boardId)
    }
    if keywordExists {
        query = query.Where("title like ?", "%"+keyword+"%")
    }
    var postCount int
    if err := config.DB.Table("(?) as a", query).Select("count(*)").Find(&postCount).Error; err != nil {
        c.Status(400)
        return
    }

    query = query.Limit(size).Offset((page - 1) * size)

    var posts []models.Post
    query.Find(&posts)

    c.IndentedJSON(200, gin.H {
        "nowPage": page,
        "pageCount": int(postCount / size) + 1,
        "pageSize": size,
        "posts": posts,
    })
}

func GetThumbnails(c *gin.Context) { 
    size, err := strconv.Atoi(c.DefaultQuery("size", "15"))
    if err != nil {
        c.Writer.WriteHeader(400)
        return
    }
    thumbnails := services.GetThumbnails(size)
    c.IndentedJSON(200, thumbnails)
}

func GetSelectedThumbnails(c *gin.Context) {
    thumbnails := services.GetSelectedThumbnails()
    c.IndentedJSON(200, thumbnails)
}

func GetPost(c *gin.Context) {
    postId, err1 := strconv.Atoi(c.Param("postId"))
    if err1 != nil {
        c.JSON(400, err1.Error())
        return
    }
    var post models.Post
    var err2 error
    token := c.Request.Header.Get("Authorization")
    if err := services.VerifyToken(token); token =="" || err != nil {
        post, err2 = services.GetEnabledPost(postId)
    } else {
        post, err2 = services.GetPost(postId)
    }
    if err2 == gorm.ErrRecordNotFound {
        c.Status(404)
    } else {
        c.IndentedJSON(200, post)
    }
}

func InsertPost(c *gin.Context) {
    requestBody := &models.Post{}
    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
    } else if err := services.InsertPost(requestBody); err != nil {
        c.JSON(400, err.Error())
    } else {
        c.Status(200)
    }
}

func UpdatePost(c *gin.Context) {
    postId, err := strconv.Atoi(c.Param("postId"))
    if err != nil {
        c.JSON(400, err.Error())
        return
    } 
    requestBody := &models.Post{}
    requestBody.PostID = postId
    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
    } else {
        err := services.UpdatePost(requestBody)
        if err == gorm.ErrRecordNotFound {
            c.Status(404)
        } else if err != nil {
            c.JSON(400, err.Error())
        } else {
            c.Status(200)
        }
    }
}

func ResetSelectedPost(c *gin.Context) {
    requestBody := &[]int{}
    if err := c.ShouldBind(requestBody); err != nil {
        c.JSON(400, err.Error())
    } else {
        if err := services.ResetSelectedPost(requestBody); err != nil {
            c.JSON(400, err.Error())
        } else {
            c.Status(200)
        }
    }
}

func DeletePost(c *gin.Context) {
    postId, err := strconv.Atoi(c.Param("postId"))
    if err != nil {
        c.JSON(400, err.Error())
        return
    } 
    
    err = services.DeletePost(postId)
    if err == gorm.ErrRecordNotFound {
        c.Status(404)
    } else if err != nil {
        c.JSON(400, err.Error())
    } else {
        c.Status(200)
    }
}
