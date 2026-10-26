package tag

// TagRepository 标签仓储接口
type TagRepository interface {
	GetHotTags(limit int) ([]map[string]interface{}, error)
	GetOrCreate(name string) (*Tag, error)
	GetByName(name string) (*Tag, error)
	UpdatePostTags(postID int64, tagIDs []int64) error
	GetPostTags(postID int64) ([]string, error)
	AddUserSkill(userID int64, tagID int64) error
	RemoveUserSkill(userID int64, tagID int64) error
	GetUserSkills(userID int64) ([]Tag, error)
	GetUserIDsByTag(tagID int64) ([]int64, error)
}
