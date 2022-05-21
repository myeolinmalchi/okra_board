package services

import (
	"errors"
	"okra_board/config"
	"okra_board/models"
	"gorm.io/gorm"
)

func GetPosts(page, size int) (posts []models.Post) {
    config.DB.Limit(size).Offset((page-1)*size).Find(&posts)
    return
}

func GetEnabledPosts(page, size int) (posts [] models.Post) {
    config.DB.
        Where("status = ?", true).
        Limit(size).Offset((page-1)*size).
        Find(&posts)
    return
}

func GetPostsWithBoardID(boardId, page, size int) (posts []models.Post) {
    config.DB.
        Where("board_id = ?", boardId).
        Limit(size).Offset((page-1)*size).
        Find(&posts)
    return
}

func GetEnabledPostsWithBoardID(boardId, page, size int) (posts []models.Post) {
    config.DB.
        Where("status = ?", true).
        Where("board_id = ?", boardId).
        Limit(size).Offset((page-1)*size).
        Find(&posts)
    return
}

func GetPost(postId int) (models.Post, error) {
    post := models.Post{}
    err := config.DB.
        Where("post_id = ?", postId).
        First(&post).
        Error
    if err != nil {
        errors.Is(err, gorm.ErrRecordNotFound)
    } else {
        config.DB.Model(&post).Association("Contents").Find(&post.Contents)
    }
    return post, err
}

func GetEnabledPost(postId int) (models.Post, error) {
    post := models.Post{}
    err := config.DB.
        Where(&models.Post{PostID:postId, Status:true}).
        First(&post).
        Error
    if err != nil {
        errors.Is(err, gorm.ErrRecordNotFound)
    } else {
        config.DB.Model(&post).Association("Contents").Find(&post.Contents)
    }
    return post, err
}

func GetThumbnails(size int) (thumbnails []models.Thumbnail) {
    config.DB.Model(&models.Post{}).Limit(size).Find(&thumbnails)
    return
}

func GetSelectedThumbnails() (thumbnails []models.Thumbnail) {
    config.DB.Table("posts").Where("selected = ?",true).Find(&thumbnails)
    return
}

func InsertPost(post *models.Post) error {
    for i := 0; i < len(post.Contents); i++ {
        post.Contents[i].Sequence = i + 1 
        if post.Contents[i].Type == "image" {

        }
    }
    return config.DB.Create(post).Error
}

func checkPostExists(postId int) (exists bool) {
    config.DB.Table("posts").
        Select("count(*) > 0").
        Where("post_id = ?", postId).
        Find(&exists)
    return 
}

func UpdatePost(post *models.Post) error {
    return config.DB.Transaction(func(tx *gorm.DB) error {
        if !checkPostExists(post.PostID) {
            return gorm.ErrRecordNotFound
        } else if err := tx.
            Delete(&models.Content{}, "post_id = ?", post.PostID).
            Error; err != nil {
            return err
        } else if err := tx.
            Session(&gorm.Session{FullSaveAssociations:true}).
            Save(&post).
            Error; err != nil {
            return err
        } else {
            return nil
        }
    })
}

func DeletePost(postId int) error {
    if !checkPostExists(postId) {
        return gorm.ErrRecordNotFound
    } else {
        return config.DB.Delete(&models.Post{}, "post_id = ?", postId).Error
    }
}

func ResetSelectedPost(ids *[]int) error {
    return config.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.
            Table("posts").
            Where("selected = ?", true).
            Update("selected", false).
            Error; err!= nil {
            return err
        } else if err := tx.
            Model(&models.Post{}).
            Where(ids).
            Update("selected", true).
            Error; err != nil {
            return err
        } else {
            return nil
        }
    })
}
