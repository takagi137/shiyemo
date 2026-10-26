package interaction

import "time"

// Like 点赞模型
type Like struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"index"`
	PostID    int64     `json:"post_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// Bookmark 收藏模型
type Bookmark struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"index"`
	PostID    int64     `json:"post_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Bookmark) TableName() string {
	return "bookmarks"
}
