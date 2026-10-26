package notification

import (
	"errors"
	"log"

	"github.com/shiyemo/server/internal/websocket"
	"gorm.io/gorm"
)

// Service 通知服务
type Service struct {
	repo RepoInterface
}

// NewService 创建通知服务
func NewService() *Service {
	return &Service{
		repo: NewRepo(),
	}
}

// NewServiceWithRepo 使用注入的 repo 创建通知服务（测试用）
func NewServiceWithRepo(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

// Create 创建通知
func (s *Service) Create(userID, actorID int64, nType string, targetID *int64, targetType *string, content string) error {
	// 不给自己发通知
	if userID == actorID {
		return nil
	}

	notification := &Notification{
		UserID:     userID,
		ActorID:    actorID,
		Type:       nType,
		TargetID:   targetID,
		TargetType: targetType,
		Content:    content,
		IsRead:     false,
	}

	err := s.repo.Create(notification)
	if err != nil {
		return err
	}

	// 通过 WebSocket 推送通知
	hub := websocket.GetHub()
	if hub != nil && hub.IsUserOnline(userID) {
		// 获取未读数量
		unreadCount, _ := s.repo.CountUnread(userID)
		hub.SendToUser(userID, "notification", map[string]interface{}{
			"id":            notification.ID,
			"type":          nType,
			"content":       content,
			"unread_count":  unreadCount,
			"created_at":    notification.CreatedAt,
		})
		log.Printf("WebSocket notification sent to user %d", userID)
	}

	return nil
}

// List 获取通知列表
func (s *Service) List(userID int64, page, pageSize int) ([]Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.ListByUserID(userID, page, pageSize)
}

// ListByType 按类型获取通知列表
func (s *Service) ListByType(userID int64, nType string, unreadOnly bool, page, pageSize int) ([]Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.ListByType(userID, nType, unreadOnly, page, pageSize)
}

// BatchMarkAsRead 批量标记已读
func (s *Service) BatchMarkAsRead(ids []int64, userID int64) error {
	if len(ids) == 0 {
		return errors.New("通知ID列表不能为空")
	}
	return s.repo.BatchMarkAsRead(ids, userID)
}

// BatchDelete 批量删除
func (s *Service) BatchDelete(ids []int64, userID int64) error {
	if len(ids) == 0 {
		return errors.New("通知ID列表不能为空")
	}
	return s.repo.BatchDelete(ids, userID)
}

// GetDetail 获取通知详情
func (s *Service) GetDetail(id int64) (*Notification, error) {
	return s.repo.GetByID(id)
}

// MarkAsRead 标记单条已读
func (s *Service) MarkAsRead(id, userID int64) error {
	notification, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return err
	}

	// 验证是否是用户自己的通知
	if notification.UserID != userID {
		return errors.New("无权操作该通知")
	}

	return s.repo.MarkAsRead(id, userID)
}

// MarkAllAsRead 标记全部已读
func (s *Service) MarkAllAsRead(userID int64) error {
	return s.repo.MarkAllAsRead(userID)
}

// GetUnreadCount 获取未读数量
func (s *Service) GetUnreadCount(userID int64) (int64, error) {
	return s.repo.CountUnread(userID)
}
