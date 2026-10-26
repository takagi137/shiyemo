package tag

import "time"

// Tag 标签
type Tag struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;size:50;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// PostTag 帖子标签关联
type PostTag struct {
	PostID int64 `gorm:"primaryKey"`
	TagID  int64 `gorm:"primaryKey"`
}

// TableName 表名
func (Tag) TableName() string      { return "tags" }
func (PostTag) TableName() string  { return "post_tags" }
