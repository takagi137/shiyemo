package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// SkillTagService 测试
// ============================================================

func TestSkillTagService_AddUserSkill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		skill, err := service.AddUserSkill(1, "Go")

		require.NoError(t, err)
		assert.NotNil(t, skill)
		assert.Equal(t, "Go", skill.Name)
	})

	t.Run("existing_tag", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		// 第一次添加
		skill1, _ := service.AddUserSkill(1, "Go")
		// 第二次添加同名标签（复用）
		skill2, err := service.AddUserSkill(1, "Go")

		require.NoError(t, err)
		assert.Equal(t, skill1.ID, skill2.ID)
	})

	t.Run("get_or_create_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetCreateError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		_, err := service.AddUserSkill(1, "Go")

		require.Error(t, err)
	})

	t.Run("add_skill_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetCreateError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		_, err := service.AddUserSkill(1, "Go")

		require.Error(t, err)
	})
}

func TestSkillTagService_RemoveUserSkill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		// 先添加
		service.AddUserSkill(1, "Go")

		// 移除
		err := service.RemoveUserSkill(1, "Go")

		require.NoError(t, err)
	})

	t.Run("tag_not_found", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		err := service.RemoveUserSkill(1, "NotExist")

		require.Error(t, err)
	})

	t.Run("get_by_name_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		err := service.RemoveUserSkill(1, "Go")

		require.Error(t, err)
	})

	t.Run("remove_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetDeleteError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		// 先创建标签
		repo.GetOrCreate("Go")

		err := service.RemoveUserSkill(1, "Go")

		require.Error(t, err)
	})
}

func TestSkillTagService_GetUserSkills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		// 先添加技能
		service.AddUserSkill(1, "Go")
		service.AddUserSkill(1, "Python")

		// 获取技能
		skills, err := service.GetUserSkills(1)

		require.NoError(t, err)
		assert.Len(t, skills, 2)
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		skills, err := service.GetUserSkills(999)

		require.NoError(t, err)
		assert.Len(t, skills, 0)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		_, err := service.GetUserSkills(1)

		require.Error(t, err)
	})
}

func TestSkillTagService_BatchAddUserSkills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		skills, err := service.BatchAddUserSkills(1, []string{"Go", "Python", "Java"})

		require.NoError(t, err)
		assert.Len(t, skills, 3)
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		skills, err := service.BatchAddUserSkills(1, []string{})

		require.NoError(t, err)
		assert.Len(t, skills, 0)
	})

	t.Run("error_on_second", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		// 第一个成功，第二个失败
		repo.SetCreateError(assert.AnError)
		_, err := service.BatchAddUserSkills(1, []string{"Go"})

		require.Error(t, err)
	})
}

func TestSkillTagService_SearchBySkill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		// 添加用户技能
		service.AddUserSkill(1, "Go")
		service.AddUserSkill(2, "Go")

		// 搜索
		userIDs, err := service.SearchBySkill("Go")

		require.NoError(t, err)
		assert.Len(t, userIDs, 2)
		assert.Contains(t, userIDs, int64(1))
		assert.Contains(t, userIDs, int64(2))
	})

	t.Run("tag_not_found", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewSkillTagServiceWithDeps(repo)

		userIDs, err := service.SearchBySkill("NotExist")

		require.NoError(t, err)
		assert.Len(t, userIDs, 0)
	})

	t.Run("get_by_name_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		_, err := service.SearchBySkill("Go")

		require.Error(t, err)
	})

	t.Run("get_user_ids_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.GetOrCreate("Go")
		repo.SetGetError(assert.AnError)
		service := NewSkillTagServiceWithDeps(repo)

		_, err := service.SearchBySkill("Go")

		require.Error(t, err)
	})
}

func TestUserSkillTag_TableName(t *testing.T) {
	tag := UserSkillTag{}
	assert.Equal(t, "user_skill_tags", tag.TableName())
}

func TestSkillTag_Struct(t *testing.T) {
	tag := SkillTag{ID: 1, Name: "Go"}
	assert.Equal(t, int64(1), tag.ID)
	assert.Equal(t, "Go", tag.Name)
}
