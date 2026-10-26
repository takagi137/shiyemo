package comment

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	PostID         int64          `json:"post_id" gorm:"index"`
	UserID         int64          `json:"user_id" gorm:"index"`
	ParentID       *int64         `json:"parent_id" gorm:"index"` // NULL 表示主评论
	Content        string         `json:"content" gorm:"type:text"`
	LikeCount      int            `json:"like_count" gorm:"default:0"`
	IsAnonymous    bool           `json:"is_anonymous" gorm:"default:false"` // 是否匿名
	AnonymousName string         `json:"anonymous_name" gorm:"type:varchar(50)"` // 匿名名称
	Status         int            `json:"status" gorm:"default:1"` // 0=草稿 1=已发布 2=隐藏 3=删除
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// IsRoot 是否是主评论
func (c *Comment) IsRoot() bool {
	return c.ParentID == nil
}

// IsDeleted 是否已删除
func (c *Comment) IsDeleted() bool {
	return c.Status == 3 || !c.DeletedAt.Time.IsZero()
}
