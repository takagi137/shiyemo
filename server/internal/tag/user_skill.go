package tag

import "time"

// UserSkillTag 用户技能标签关联
type UserSkillTag struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	TagID     int64     `json:"tag_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (UserSkillTag) TableName() string {
	return "user_skill_tags"
}

// SkillTag 技能标签响应
type SkillTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
