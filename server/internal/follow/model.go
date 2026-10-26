package follow

import "time"

// Follow 关注关系模型
type Follow struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	FollowerID  int64     `json:"follower_id" gorm:"index"`
	FollowingID int64     `json:"following_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (Follow) TableName() string {
	return "follows"
}
