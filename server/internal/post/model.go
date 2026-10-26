package post

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Post 帖子模型
type Post struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64          `json:"user_id" gorm:"index"`
	Type           string         `json:"type" gorm:"type:varchar(20);index"` // emotion/insight/case/action/achievement
	Title          string         `json:"title" gorm:"type:varchar(200)"`
	Content        string         `json:"content" gorm:"type:text"`
	Images         string         `json:"images" gorm:"type:jsonb"` // JSON 数组
	Tags           string         `json:"tags" gorm:"type:jsonb"`   // JSON 数组
	LikeCount      int            `json:"like_count" gorm:"default:0"`
	CommentCount   int            `json:"comment_count" gorm:"default:0"`
	IsTop          bool           `json:"is_top" gorm:"default:false"`
	IsFeatured     bool           `json:"is_featured" gorm:"default:false"`
	IsAnonymous    bool           `json:"is_anonymous" gorm:"default:false"` // 是否匿名
	AnonymousName string         `json:"anonymous_name" gorm:"type:varchar(50)"` // 匿名名称
	Status         int            `json:"status" gorm:"default:1"` // 0=草稿 1=已发布 2=隐藏 3=删除
	ReviewStatus   string         `json:"review_status" gorm:"type:varchar(20);default:'approved';index"` // pending/approved/rejected
	ReviewedAt     *time.Time     `json:"reviewed_at"`
	ReviewedBy     *int64         `json:"reviewed_by"`
	RejectReason   string         `json:"reject_reason" gorm:"type:text"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// GetImages 获取图片列表
func (p *Post) GetImages() []string {
	var images []string
	if p.Images != "" {
		json.Unmarshal([]byte(p.Images), &images)
	}
	if images == nil {
		images = []string{}
	}
	return images
}

// GetTags 获取标签列表
func (p *Post) GetTags() []string {
	var tags []string
	if p.Tags != "" {
		json.Unmarshal([]byte(p.Tags), &tags)
	}
	if tags == nil {
		tags = []string{}
	}
	return tags
}

// SetImages 设置图片列表
func (p *Post) SetImages(images []string) {
	if images == nil {
		images = []string{}
	}
	data, _ := json.Marshal(images)
	p.Images = string(data)
}

// SetTags 设置标签列表
func (p *Post) SetTags(tags []string) {
	if tags == nil {
		tags = []string{}
	}
	data, _ := json.Marshal(tags)
	p.Tags = string(data)
}

// IsPublished 是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == 1 && p.DeletedAt.Time.IsZero()
}
