package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId uint	`json:"user_id" gorm:"not null"`
	CategoryId uint `json:"category_id" gorm:"not null"`
	Category *Category
	Title string `json:"title" gorm:"type:varchar(50);not null"`
	HeadImg string `json:"head_img"`
	Content string `json:"content" gorm:"type:text;not null"`
	CreatedAt Time `json:"created-at" gorm:"type:timestamp"`
	UpdateAt Time `json:"update_at" gorm:"type:timestamp"`
}

func (post *Post) BeforeCreate(scope *gorm.DB) error {
	post.ID = uuid.NewV4()
	return nil
}
