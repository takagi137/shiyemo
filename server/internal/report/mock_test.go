package report

import (
	"errors"
)

// MockReportRepo 举报仓储Mock实现
type MockReportRepo struct {
	reports []Report
	nextID  int64
	err     error
	exists  bool
}

// NewMockReportRepo 创建Mock仓储
func NewMockReportRepo() *MockReportRepo {
	return &MockReportRepo{
		reports: make([]Report, 0),
		nextID:  1,
	}
}

// SetError 设置返回错误
func (m *MockReportRepo) SetError(err error) {
	m.err = err
}

// SetExists 设置已存在
func (m *MockReportRepo) SetExists(exists bool) {
	m.exists = exists
}

// Create 创建举报
func (m *MockReportRepo) Create(report *Report) error {
	if m.err != nil {
		return m.err
	}
	report.ID = m.nextID
	m.nextID++
	m.reports = append(m.reports, *report)
	return nil
}

// ExistsByUserAndTarget 检查用户是否已举报过该目标
func (m *MockReportRepo) ExistsByUserAndTarget(userID int64, targetType string, targetID int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.exists, nil
}

// List 获取举报列表
func (m *MockReportRepo) List(page, size int) ([]Report, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	total := int64(len(m.reports))
	start := (page - 1) * size
	end := start + size
	if start > len(m.reports) {
		return []Report{}, total, nil
	}
	if end > len(m.reports) {
		end = len(m.reports)
	}
	return m.reports[start:end], total, nil
}

// UpdateStatus 更新举报状态
func (m *MockReportRepo) UpdateStatus(id int64, status int) error {
	if m.err != nil {
		return m.err
	}
	for i, report := range m.reports {
		if report.ID == id {
			m.reports[i].Status = status
			return nil
		}
	}
	return errors.New("report not found")
}
