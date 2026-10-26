package export

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock 实现
// ============================================================

type mockUserProvider struct {
	users    []auth.User
	count    int64
	listErr  error
	countErr error
}

func (m *mockUserProvider) GetAll() ([]auth.User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.users, nil
}

func (m *mockUserProvider) Count() (int64, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return m.count, nil
}

type mockPostProvider struct {
	posts    []post.Post
	total    int64
	listErr  error
	countErr error
}

func (m *mockPostProvider) List(postType string, page, size int) ([]post.Post, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	if postType == "" {
		return m.posts, m.total, nil
	}
	var filtered []post.Post
	for _, p := range m.posts {
		if p.Type == postType {
			filtered = append(filtered, p)
		}
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockPostProvider) Count() (int64, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return m.total, nil
}

// ============================================================
// ExportUsers 测试
// ============================================================

func TestExportUsers_Success(t *testing.T) {
	email1 := "test1@example.com"
	email2 := "test2@example.com"
	
	mockUsers := []auth.User{
		{
			ID:         1,
			Email:      &email1,
			Nickname:   "用户1",
			Avatar:     "avatar1.jpg",
			Bio:        "简介1",
			Company:    "公司1",
			Occupation: "职业1",
			IsAdmin:    false,
			CreatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:         2,
			Email:      &email2,
			Nickname:   "用户2",
			IsAdmin:    true,
			CreatedAt:  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	mock := &mockUserProvider{users: mockUsers, count: 2}
	svc := NewExportServiceWithDeps(mock, nil)

	data, err := svc.ExportUsers()

	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// 解析 CSV
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 验证表头 + 2条数据
	assert.Len(t, records, 3)
	assert.Equal(t, []string{"ID", "邮箱", "昵称", "头像", "简介", "公司", "职业", "是否管理员", "注册时间"}, records[0])
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "test1@example.com", records[1][1])
	assert.Equal(t, "用户1", records[1][2])
	assert.Equal(t, "false", records[1][7])
	assert.Equal(t, "2", records[2][0])
	assert.Equal(t, "true", records[2][7])
}

func TestExportUsers_EmptyList(t *testing.T) {
	mock := &mockUserProvider{users: []auth.User{}, count: 0}
	svc := NewExportServiceWithDeps(mock, nil)

	data, err := svc.ExportUsers()

	require.NoError(t, err)
	
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 只有表头
	assert.Len(t, records, 1)
}

func TestExportUsers_Error(t *testing.T) {
	mock := &mockUserProvider{listErr: assert.AnError}
	svc := NewExportServiceWithDeps(mock, nil)

	data, err := svc.ExportUsers()

	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestExportUsers_NilEmail(t *testing.T) {
	mockUsers := []auth.User{
		{
			ID:        1,
			Email:     nil,
			Nickname:  "用户1",
			CreatedAt: time.Now(),
		},
	}

	mock := &mockUserProvider{users: mockUsers, count: 1}
	svc := NewExportServiceWithDeps(mock, nil)

	data, err := svc.ExportUsers()

	require.NoError(t, err)
	
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 邮箱字段应该为空字符串
	assert.Equal(t, "", records[1][1])
}

// ============================================================
// ExportPosts 测试
// ============================================================

func TestExportPosts_Success(t *testing.T) {
	mockPosts := []post.Post{
		{
			ID:           1,
			UserID:       100,
			Type:         "emotion",
			Title:        "标题1",
			Content:      "内容1",
			LikeCount:    10,
			CommentCount: 5,
			Status:       1,
			CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:           2,
			UserID:       200,
			Type:         "insight",
			Title:        "标题2",
			Content:      "内容2",
			LikeCount:    20,
			CommentCount: 15,
			Status:       1,
			CreatedAt:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	mock := &mockPostProvider{posts: mockPosts, total: 2}
	svc := NewExportServiceWithDeps(nil, mock)

	data, err := svc.ExportPosts("")

	require.NoError(t, err)
	assert.NotEmpty(t, data)

	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 验证表头 + 2条数据
	assert.Len(t, records, 3)
	assert.Equal(t, []string{"ID", "用户ID", "类型", "标题", "内容", "点赞数", "评论数", "状态", "创建时间"}, records[0])
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "100", records[1][1])
	assert.Equal(t, "emotion", records[1][2])
	assert.Equal(t, "10", records[1][5])
}

func TestExportPosts_FilterByType(t *testing.T) {
	mockPosts := []post.Post{
		{ID: 1, Type: "emotion", Title: "情感", Status: 1, CreatedAt: time.Now()},
		{ID: 2, Type: "insight", Title: "洞察", Status: 1, CreatedAt: time.Now()},
		{ID: 3, Type: "emotion", Title: "情感2", Status: 1, CreatedAt: time.Now()},
	}

	mock := &mockPostProvider{posts: mockPosts, total: 3}
	svc := NewExportServiceWithDeps(nil, mock)

	data, err := svc.ExportPosts("emotion")

	require.NoError(t, err)

	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 只有 emotion 类型的帖子
	assert.Len(t, records, 3) // 表头 + 2条
}

func TestExportPosts_Error(t *testing.T) {
	mock := &mockPostProvider{listErr: assert.AnError}
	svc := NewExportServiceWithDeps(nil, mock)

	data, err := svc.ExportPosts("")

	assert.Error(t, err)
	assert.Nil(t, data)
}

// ============================================================
// ExportStats 测试
// ============================================================

func TestExportStats_Success(t *testing.T) {
	mockUsers := &mockUserProvider{count: 100}
	mockPosts := &mockPostProvider{total: 500}
	svc := NewExportServiceWithDeps(mockUsers, mockPosts)

	stats, err := svc.ExportStats()

	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(100), stats["user_count"])
	assert.Equal(t, int64(500), stats["post_count"])
	assert.NotEmpty(t, stats["exported_at"])
}

func TestExportStats_UserCountError(t *testing.T) {
	mockUsers := &mockUserProvider{countErr: assert.AnError}
	mockPosts := &mockPostProvider{total: 500}
	svc := NewExportServiceWithDeps(mockUsers, mockPosts)

	stats, err := svc.ExportStats()

	assert.Error(t, err)
	assert.Nil(t, stats)
}

func TestExportStats_PostCountError(t *testing.T) {
	mockUsers := &mockUserProvider{count: 100}
	mockPosts := &mockPostProvider{countErr: assert.AnError}
	svc := NewExportServiceWithDeps(mockUsers, mockPosts)

	stats, err := svc.ExportStats()

	assert.Error(t, err)
	assert.Nil(t, stats)
}
