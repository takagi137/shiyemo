package feedback

import "errors"

// RepoInterface 仓储接口（支持 mock 测试）
type RepoInterface interface {
	Create(feedback *Feedback) error
	List(page, pageSize int, status *int, feedbackType string) ([]FeedbackResponse, int64, error)
	ListByUser(userID int64, page, pageSize int) ([]Feedback, int64, error)
	UpdateStatus(id int64, status int8, reply string) error
}

// Service 反馈服务
type Service struct {
	repo RepoInterface
}

// NewService 创建服务
func NewService() *Service {
	return &Service{repo: NewRepo()}
}

// NewServiceWithDeps 创建服务（注入依赖，用于测试）
func NewServiceWithDeps(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

// Create 创建反馈
func (s *Service) Create(userID int64, feedbackType, title, content, contact string) error {
	if title == "" {
		return errors.New("标题不能为空")
	}
	if content == "" {
		return errors.New("内容不能为空")
	}

	// 验证类型
	validTypes := map[string]bool{"bug": true, "suggestion": true, "other": true}
	if !validTypes[feedbackType] {
		feedbackType = "other"
	}

	feedback := &Feedback{
		UserID:  userID,
		Type:    feedbackType,
		Title:   title,
		Content: content,
		Contact: contact,
		Status:  0,
	}

	return s.repo.Create(feedback)
}

// List 获取反馈列表（管理员）
func (s *Service) List(page, pageSize int, status *int, feedbackType string) ([]FeedbackResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.List(page, pageSize, status, feedbackType)
}

// ListByUser 获取用户的反馈列表
func (s *Service) ListByUser(userID int64, page, pageSize int) ([]Feedback, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListByUser(userID, page, pageSize)
}

// UpdateStatus 更新反馈状态（管理员）
func (s *Service) UpdateStatus(id int64, status int8, reply string) error {
	// 验证状态
	if status != 0 && status != 1 && status != 2 {
		return errors.New("无效的状态值")
	}
	return s.repo.UpdateStatus(id, status, reply)
}
