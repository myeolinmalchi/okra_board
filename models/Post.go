package models


import (
    "time"
)

type Post struct {
    PostID      int         `json:"postId,omitemtpy" gorm:"primaryKey"`
    BoardID     int         `json:"boardId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
    AddedDate   time.Time   `json:"addedDate,omitempty" gorm:"->"`
    Status      bool        `json:"status"`
    Selected    bool        `json:"selected"`
    Contents    []Content   `json:"contents,omitempty" gorm:"foreignKey:PostID"`
}

// Response Only
type Thumbnail struct {
    PostID      int         `json:"postId"`
    Title       string      `json:"title"`
    Thumbnail   string      `json:"thumbnail"`
}

type Content struct {
    PostID      int         `json:"postId,omitemtpy"`
    ContentID   int         `json:"contentId,omitempty" gorm:"primaryKey"`
    Type        string      `json:"type"`
    Content     string      `json:"content"`
    Sequence    int         `json:"sequence,omitempty"`
}
