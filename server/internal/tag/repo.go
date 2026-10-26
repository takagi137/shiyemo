package tag

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// TagRepo 标签仓库
type TagRepo struct {
	db *gorm.DB
}

// NewTagRepo 创建标签仓库
func NewTagRepo() *TagRepo {
	return &TagRepo{db: database.GetDB()}
}

// NewTagRepoWithDB 使用注入的 DB 创建标签仓库（测试用）
func NewTagRepoWithDB(db *gorm.DB) *TagRepo {
	return &TagRepo{db: db}
}

// GetByName 根据名称获取标签
func (r *TagRepo) GetByName(name string) (*Tag, error) {
	var tag Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// GetOrCreate 获取或创建标签
func (r *TagRepo) GetOrCreate(name string) (*Tag, error) {
	var tag Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		tag = Tag{Name: name}
		err = r.db.Create(&tag).Error
	}
	return &tag, err
}

// GetHotTags 获取热门标签
func (r *TagRepo) GetHotTags(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := r.db.Raw(`
		SELECT t.id, t.name, COUNT(pt.post_id) as count
		FROM tags t
		LEFT JOIN post_tags pt ON t.id = pt.tag_id
		GROUP BY t.id, t.name
		ORDER BY count DESC
		LIMIT ?
	`, limit).Scan(&results).Error
	
	return results, err
}

// GetPostTags 获取帖子的标签
func (r *TagRepo) GetPostTags(postID int64) ([]string, error) {
	var names []string
	err := r.db.Raw(`
		SELECT t.name
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = ?
	`, postID).Scan(&names).Error
	return names, err
}

// UpdatePostTags 更新帖子标签
func (r *TagRepo) UpdatePostTags(postID int64, tagIDs []int64) error {
	// 删除旧关联
	r.db.Where("post_id = ?", postID).Delete(&PostTag{})
	
	// 插入新关联
	if len(tagIDs) > 0 {
		var postTags []PostTag
		for _, tagID := range tagIDs {
			postTags = append(postTags, PostTag{PostID: postID, TagID: tagID})
		}
		return r.db.Create(&postTags).Error
	}
	return nil
}

// AddUserSkill 添加用户技能标签
func (r *TagRepo) AddUserSkill(userID, tagID int64) error {
	// 检查是否已存在
	var count int64
	r.db.Model(&UserSkillTag{}).Where("user_id = ? AND tag_id = ?", userID, tagID).Count(&count)
	if count > 0 {
		return nil
	}

	return r.db.Create(&UserSkillTag{
		UserID: userID,
		TagID:  tagID,
	}).Error
}

// RemoveUserSkill 移除用户技能标签
func (r *TagRepo) RemoveUserSkill(userID, tagID int64) error {
	return r.db.Where("user_id = ? AND tag_id = ?", userID, tagID).Delete(&UserSkillTag{}).Error
}

// GetUserSkills 获取用户的技能标签
func (r *TagRepo) GetUserSkills(userID int64) ([]Tag, error) {
	var tags []Tag
	err := r.db.Raw(`
		SELECT t.id, t.name
		FROM tags t
		JOIN user_skill_tags ust ON t.id = ust.tag_id
		WHERE ust.user_id = ?
		ORDER BY ust.created_at DESC
	`, userID).Scan(&tags).Error
	return tags, err
}

// GetUserIDsByTag 根据标签获取用户ID列表
func (r *TagRepo) GetUserIDsByTag(tagID int64) ([]int64, error) {
	var userIDs []int64
	err := r.db.Model(&UserSkillTag{}).Where("tag_id = ?", tagID).Pluck("user_id", &userIDs).Error
	return userIDs, err
}
