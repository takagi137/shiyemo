package report

// ReportRepository 举报仓储接口
type ReportRepository interface {
	Create(report *Report) error
	ExistsByUserAndTarget(userID int64, targetType string, targetID int64) (bool, error)
	List(page, size int) ([]Report, int64, error)
	UpdateStatus(id int64, status int) error
}
