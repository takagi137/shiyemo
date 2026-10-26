package post

import (
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/moderation"
	"gorm.io/gorm"
)

// ============================================================
// Mock PostRepo
// ============================================================

type mockPostRepo struct {
	posts  map[int64]*Post
	nextID int64
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{posts: make(map[int64]*Post)}
}

func (m *mockPostRepo) Create(post *Post) error {
	m.nextID++
	post.ID = m.nextID
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) GetByID(id int64) (*Post, error) {
	if p, ok := m.posts[id]; ok {
		return p, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockPostRepo) List(postType string, page, size int) ([]Post, int64, error) {
	var result []Post
	for _, p := range m.posts {
		if p.Status == 1 && (postType == "" || p.Type == postType) {
			result = append(result, *p)
		}
	}
	total := int64(len(result))

	start := (page - 1) * size
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + size
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (m *mockPostRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	if p, ok := m.posts[id]; ok {
		if v, ok := fields["title"]; ok {
			p.Title = v.(string)
		}
		if v, ok := fields["content"]; ok {
			p.Content = v.(string)
		}
		if v, ok := fields["type"]; ok {
			p.Type = v.(string)
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockPostRepo) Delete(id int64) error {
	if _, ok := m.posts[id]; ok {
		delete(m.posts, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockPostRepo) Count() (int64, error) {
	var count int64
	for _, p := range m.posts {
		if p.Status == 1 {
			count++
		}
	}
	return count, nil
}

func (m *mockPostRepo) GetByUserIDs(userIDs []int64, page, size int) ([]Post, int64, error) {
	idSet := make(map[int64]bool)
	for _, id := range userIDs {
		idSet[id] = true
	}
	var result []Post
	for _, p := range m.posts {
		if p.Status == 1 && idSet[p.UserID] {
			result = append(result, *p)
		}
	}
	total := int64(len(result))

	start := (page - 1) * size
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + size
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (m *mockPostRepo) GetFollowingIDs(userID int64) ([]int64, error) {
	return nil, nil
}

func (m *mockPostRepo) ListByCohort(postType, city, industry string, page, size int) ([]Post, int64, error) {
	var result []Post
	for _, p := range m.posts {
		if p.Status == 1 && (postType == "" || p.Type == postType) {
			result = append(result, *p)
		}
	}
	total := int64(len(result))

	start := (page - 1) * size
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + size
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

// ============================================================
// Mock PostRepo with errors
// ============================================================

type errorPostRepo struct {
	*mockPostRepo
	createErr      error
	getErr         error
	getFollowingErr error
}

func (m *errorPostRepo) Create(post *Post) error {
	if m.createErr != nil {
		return m.createErr
	}
	return m.mockPostRepo.Create(post)
}

func (m *errorPostRepo) GetByID(id int64) (*Post, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.mockPostRepo.GetByID(id)
}

func (m *errorPostRepo) GetFollowingIDs(userID int64) ([]int64, error) {
	if m.getFollowingErr != nil {
		return nil, m.getFollowingErr
	}
	return m.mockPostRepo.GetFollowingIDs(userID)
}

// ============================================================
// Mock PostIndexer
// ============================================================

type mockPostIndexer struct {
	indexErr  error
	deleteErr error
	indexed   []*Post
	deleted   []int64
}

func newMockPostIndexer() *mockPostIndexer {
	return &mockPostIndexer{}
}

func (m *mockPostIndexer) IndexPost(p *Post) error {
	if m.indexErr != nil {
		return m.indexErr
	}
	m.indexed = append(m.indexed, p)
	return nil
}

func (m *mockPostIndexer) DeletePost(postID int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deleted = append(m.deleted, postID)
	return nil
}

// ============================================================
// Mock UserGetter
// ============================================================

type mockUserGetter struct {
	users map[int64]*auth.User
}

func newMockUserGetter() *mockUserGetter {
	return &mockUserGetter{users: make(map[int64]*auth.User)}
}

func (m *mockUserGetter) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// ============================================================
// Mock ContentModerator
// ============================================================

type mockModerator struct {
	checkErr *moderation.CheckResult
}

func newMockModerator() *mockModerator {
	return &mockModerator{
		checkErr: &moderation.CheckResult{Passed: true},
	}
}

func (m *mockModerator) CheckTitle(title string) *moderation.CheckResult {
	return m.checkErr
}

func (m *mockModerator) CheckContent(text string) *moderation.CheckResult {
	return m.checkErr
}

// ============================================================
// 工具函数
// ============================================================

func createTestPost(userID int64, title, content, postType string) *Post {
	p := &Post{
		UserID:  userID,
		Type:    postType,
		Title:   title,
		Content: content,
		Status:  1,
	}
	p.SetImages(nil)
	p.SetTags(nil)
	return p
}

func seedPosts(repo *mockPostRepo, userID int64, count int) {
	for i := 0; i < count; i++ {
		repo.Create(createTestPost(userID, "帖"+string(rune('0'+i)), "内容", "emotion"))
	}
}

func createService() (*PostService, *mockPostRepo, *mockUserGetter, *mockModerator) {
	postRepo := newMockPostRepo()
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	svc := NewPostServiceWithDeps(postRepo, userRepo, moderator)
	return svc, postRepo, userRepo, moderator
}
