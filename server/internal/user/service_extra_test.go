package user

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// 错误路径 Mock
// ============================================================

type errorUserRepo struct {
	*mockUserRepo
	err error
}

func (m *errorUserRepo) GetByID(id int64) (*auth.User, error) {
	return nil, m.err
}

type errorUpdateUserRepo struct {
	*mockUserRepo
	updateFieldsErr error
}

func (m *errorUpdateUserRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return m.updateFieldsErr
}

// ============================================================
// 错误路径测试
// ============================================================

func TestUserService_GetProfile_DBError(t *testing.T) {
	mock := &errorUserRepo{mockUserRepo: newMockUserRepo(), err: errors.New("database connection refused")}
	svc := NewUserServiceWithRepo(mock)

	_, err := svc.GetProfile(100)
	require.Error(t, err)
	assert.Equal(t, bizerr.ErrDB, err)
}

func TestUserService_GetPublicProfile_DBError(t *testing.T) {
	mock := &errorUserRepo{mockUserRepo: newMockUserRepo(), err: errors.New("database connection refused")}
	svc := NewUserServiceWithRepo(mock)

	_, err := svc.GetPublicProfile(100)
	require.Error(t, err)
}

func TestUserService_GetUserStats_DBError(t *testing.T) {
	mock := &errorUserRepo{mockUserRepo: newMockUserRepo(), err: errors.New("database connection refused")}
	svc := NewUserServiceWithRepo(mock)

	_, err := svc.GetUserStats(100)
	require.Error(t, err)
}

func TestUserService_UpdateProfile_RepoError(t *testing.T) {
	mock := &errorUserRepo{mockUserRepo: newMockUserRepo(), err: errors.New("连接失败")}
	svc := NewUserServiceWithRepo(mock)

	err := svc.UpdateProfile(100, UpdateProfileRequest{Nickname: "新名"})
	require.Error(t, err)
	assert.Equal(t, bizerr.ErrDB, err)
}

func TestUserService_UpdateProfile_UpdateFieldsError(t *testing.T) {
	mock := &errorUpdateUserRepo{
		mockUserRepo:    newMockUserRepo(),
		updateFieldsErr: errors.New("更新失败"),
	}
	mock.users[100] = &auth.User{ID: 100}
	svc := NewUserServiceWithRepo(mock)

	err := svc.UpdateProfile(100, UpdateProfileRequest{Nickname: "新名"})
	require.Error(t, err)
}

// ============================================================
// GetUserPosts + UploadAvatar 测试
// ============================================================

func TestUserService_GetUserPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockUserRepoWithPosts{
			mockUserRepo: newMockUserRepo(),
			posts: map[int64][]map[string]interface{}{
				100: {{"id": 1, "title": "帖子1"}, {"id": 2, "title": "帖子2"}},
			},
		}
		svc := NewUserServiceWithRepo(mock)

		posts, total, err := svc.GetUserPosts(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("empty", func(t *testing.T) {
		svc, _ := createMockService()
		posts, total, err := svc.GetUserPosts(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})

	t.Run("pagination", func(t *testing.T) {
		postsData := make([]map[string]interface{}, 5)
		for i := range postsData {
			postsData[i] = map[string]interface{}{"id": i + 1}
		}
		mock := &mockUserRepoWithPosts{
			mockUserRepo: newMockUserRepo(),
			posts:        map[int64][]map[string]interface{}{100: postsData},
		}
		svc := NewUserServiceWithRepo(mock)

		p1, total, _ := svc.GetUserPosts(100, 1, 3)
		assert.Equal(t, int64(5), total)
		assert.Len(t, p1, 3)

		p2, _, _ := svc.GetUserPosts(100, 2, 3)
		assert.Len(t, p2, 2)
	})

	t.Run("default_params", func(t *testing.T) {
		svc, _ := createMockService()
		_, _, err := svc.GetUserPosts(100, 0, 0)
		require.NoError(t, err)
	})
}

func TestUserService_UploadAvatar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100}

		os.MkdirAll("uploads/avatars", 0755)
		defer os.RemoveAll("uploads")

		reader := bytes.NewReader([]byte("fake-image-data"))
		url, err := svc.UploadAvatar(100, reader, "avatar.png")
		require.NoError(t, err)
		assert.True(t, strings.Contains(url, "100.png"))
		assert.Equal(t, url, mock.users[100].Avatar)
	})

	t.Run("invalid_extension", func(t *testing.T) {
		svc, _ := createMockService()
		reader := bytes.NewReader([]byte("data"))
		_, err := svc.UploadAvatar(100, reader, "avatar.exe")
		require.Error(t, err)
	})
}

// ============================================================
// 补充覆盖
// ============================================================

func TestUserService_GetProfile_Cache(t *testing.T) {
	svc, mock := createMockService()
	mock.users[100] = &auth.User{ID: 100, Nickname: "缓存测试"}

	user1, err := svc.GetProfile(100)
	require.NoError(t, err)
	assert.Equal(t, "缓存测试", user1.Nickname)

	user2, err := svc.GetProfile(100)
	require.NoError(t, err)
	assert.Equal(t, "缓存测试", user2.Nickname)
}

func TestUserService_GetPublicProfile_Cache(t *testing.T) {
	svc, mock := createMockService()
	mock.users[100] = &auth.User{ID: 100, Nickname: "公开测试"}

	profile1, err := svc.GetPublicProfile(100)
	require.NoError(t, err)
	assert.Equal(t, "公开测试", profile1["nickname"])

	profile2, err := svc.GetPublicProfile(100)
	require.NoError(t, err)
	assert.Equal(t, "公开测试", profile2["nickname"])
}

func TestUserService_GetUserStats_AllFields(t *testing.T) {
	svc, mock := createMockService()
	mock.users[100] = &auth.User{ID: 100, FollowerCount: 20, FollowingCount: 12}

	stats, err := svc.GetUserStats(100)
	require.NoError(t, err)
	assert.Contains(t, stats, "post_count")
	assert.Contains(t, stats, "like_received")
	assert.Contains(t, stats, "follower_count")
	assert.Contains(t, stats, "following_count")
}

func TestUserService_UpdateProfile_AllFieldsSet(t *testing.T) {
	svc, mock := createMockService()
	mock.users[100] = &auth.User{ID: 100}

	err := svc.UpdateProfile(100, UpdateProfileRequest{
		Nickname:   "Full Update",
		Bio:        "New Bio",
		Company:    "New Company",
		Occupation: "Developer",
	})
	require.NoError(t, err)

	u := mock.users[100]
	assert.Equal(t, "Full Update", u.Nickname)
	assert.Equal(t, "New Bio", u.Bio)
	assert.Equal(t, "New Company", u.Company)
	assert.Equal(t, "Developer", u.Occupation)
}

func TestUserService_GetUserPosts_OutOfRange(t *testing.T) {
	mock := &mockUserRepoWithPosts{
		mockUserRepo: newMockUserRepo(),
		posts:        map[int64][]map[string]interface{}{100: {{"id": 1}}},
	}
	svc := NewUserServiceWithRepo(mock)

	posts, total, err := svc.GetUserPosts(100, 100, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Nil(t, posts)
}

func TestUserService_InvalidExtension(t *testing.T) {
	svc, _ := createMockService()

	exts := []string{".exe", ".pdf", ".txt", ".bmp"}
	for _, ext := range exts {
		reader := bytes.NewReader([]byte("data"))
		_, err := svc.UploadAvatar(100, reader, "avatar"+ext)
		require.Error(t, err, "扩展名 %s 应被拒绝", ext)
	}
}
