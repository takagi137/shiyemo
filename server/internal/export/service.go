package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
)

// ============================================================
// 接口定义
// ============================================================

// UserProvider 用户数据提供者
type UserProvider interface {
	GetAll() ([]auth.User, error)
	Count() (int64, error)
}

// PostProvider 帖子数据提供者
type PostProvider interface {
	List(postType string, page, size int) ([]post.Post, int64, error)
	Count() (int64, error)
}

// ============================================================
// 服务实现
// ============================================================

// ExportService 数据导出服务
type ExportService struct {
	userRepo UserProvider
	postRepo PostProvider
}

// NewExportService 创建导出服务（使用数据库）
func NewExportService() *ExportService {
	return &ExportService{
		userRepo: auth.NewAuthRepo(),
		postRepo: post.NewPostRepo(),
	}
}

// NewExportServiceWithDeps 创建导出服务（注入依赖，用于测试）
func NewExportServiceWithDeps(userRepo UserProvider, postRepo PostProvider) *ExportService {
	return &ExportService{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

// ExportUsers 导出用户数据
func (s *ExportService) ExportUsers() ([]byte, error) {
	// 获取所有用户
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 创建 CSV 缓冲区
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// 写入表头
	headers := []string{"ID", "邮箱", "昵称", "头像", "简介", "公司", "职业", "是否管理员", "注册时间"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, user := range users {
		email := ""
		if user.Email != nil {
			email = *user.Email
		}
		record := []string{
			fmt.Sprintf("%d", user.ID),
			email,
			user.Nickname,
			user.Avatar,
			user.Bio,
			user.Company,
			user.Occupation,
			fmt.Sprintf("%t", user.IsAdmin),
			user.CreatedAt.Format(time.DateTime),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// ExportPosts 导出帖子数据
func (s *ExportService) ExportPosts(postType string) ([]byte, error) {
	// 获取所有帖子
	posts, _, err := s.postRepo.List(postType, 1, 10000)
	if err != nil {
		return nil, err
	}

	// 创建 CSV 缓冲区
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// 写入表头
	headers := []string{"ID", "用户ID", "类型", "标题", "内容", "点赞数", "评论数", "状态", "创建时间"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, p := range posts {
		record := []string{
			fmt.Sprintf("%d", p.ID),
			fmt.Sprintf("%d", p.UserID),
			p.Type,
			p.Title,
			p.Content,
			fmt.Sprintf("%d", p.LikeCount),
			fmt.Sprintf("%d", p.CommentCount),
			fmt.Sprintf("%d", p.Status),
			p.CreatedAt.Format(time.DateTime),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// ExportStats 导出统计数据
func (s *ExportService) ExportStats() (map[string]interface{}, error) {
	// 用户统计
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	// 帖子统计
	postCount, err := s.postRepo.Count()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_count":  userCount,
		"post_count":  postCount,
		"exported_at": time.Now().Format(time.DateTime),
	}, nil
}
