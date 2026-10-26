package report

import bizerr "github.com/shiyemo/server/internal/common/errors"

// ReportService 举报服务
type ReportService struct {
	repo ReportRepository
}

// NewReportService 创建举报服务
func NewReportService() *ReportService {
	return &ReportService{
		repo: NewReportRepo(),
	}
}

// NewReportServiceWithDeps 创建举报服务（注入依赖，用于测试）
func NewReportServiceWithDeps(repo ReportRepository) *ReportService {
	return &ReportService{
		repo: repo,
	}
}

// Create 创建举报
func (s *ReportService) Create(userID int64, req CreateRequest) error {
	// 检查是否已举报过
	exists, err := s.repo.ExistsByUserAndTarget(userID, req.TargetType, req.TargetID)
	if err != nil {
		return bizerr.ErrDB
	}
	if exists {
		return bizerr.New(40001, "你已经举报过该内容")
	}

	report := &Report{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
		Status:     1,
	}

	if err := s.repo.Create(report); err != nil {
		return bizerr.ErrDB
	}

	return nil
}

// List 获取举报列表（管理员）
func (s *ReportService) List(page, size int) ([]Report, int64, error) {
	return s.repo.List(page, size)
}

// Handle 处理举报（管理员）
func (s *ReportService) Handle(reportID int64, status int) error {
	return s.repo.UpdateStatus(reportID, status)
}

// CreateRequest 创建举报请求
type CreateRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=post comment"`
	TargetID   int64  `json:"target_id" binding:"required"`
	Reason     string `json:"reason" binding:"required,max=500"`
}
