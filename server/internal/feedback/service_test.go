package feedback

import (
	"errors"
	"testing"
)

// mockFeedbackRepo Mock 反馈仓储
type mockFeedbackRepo struct {
	feedbacks []FeedbackResponse
	nextID    int64
	errCreate error
	errList   error
	errUpdate error
}

func newMockFeedbackRepo() *mockFeedbackRepo {
	return &mockFeedbackRepo{
		feedbacks: make([]FeedbackResponse, 0),
		nextID:    1,
	}
}

func (m *mockFeedbackRepo) Create(feedback *Feedback) error {
	if m.errCreate != nil {
		return m.errCreate
	}
	feedback.ID = m.nextID
	m.nextID++
	m.feedbacks = append(m.feedbacks, FeedbackResponse{Feedback: *feedback})
	return nil
}

func (m *mockFeedbackRepo) List(page, pageSize int, status *int, feedbackType string) ([]FeedbackResponse, int64, error) {
	if m.errList != nil {
		return nil, 0, m.errList
	}
	total := int64(len(m.feedbacks))
	offset := (page - 1) * pageSize
	if offset >= len(m.feedbacks) {
		return []FeedbackResponse{}, total, nil
	}
	end := offset + pageSize
	if end > len(m.feedbacks) {
		end = len(m.feedbacks)
	}
	return m.feedbacks[offset:end], total, nil
}

func (m *mockFeedbackRepo) ListByUser(userID int64, page, pageSize int) ([]Feedback, int64, error) {
	if m.errList != nil {
		return nil, 0, m.errList
	}
	var result []Feedback
	for _, f := range m.feedbacks {
		if f.UserID == userID {
			result = append(result, f.Feedback)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockFeedbackRepo) UpdateStatus(id int64, status int8, reply string) error {
	if m.errUpdate != nil {
		return m.errUpdate
	}
	for i, f := range m.feedbacks {
		if f.ID == id {
			m.feedbacks[i].Status = status
			m.feedbacks[i].Reply = reply
			return nil
		}
	}
	return errors.New("record not found")
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		err := svc.Create(1, "bug", "测试标题", "测试内容", "test@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(repo.feedbacks) != 1 {
			t.Fatalf("expected 1 feedback, got %d", len(repo.feedbacks))
		}
		if repo.feedbacks[0].Title != "测试标题" {
			t.Fatalf("expected title '测试标题', got '%s'", repo.feedbacks[0].Title)
		}
	})

	t.Run("empty_title", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		err := svc.Create(1, "bug", "", "内容", "")
		if err == nil {
			t.Fatal("expected error for empty title")
		}
	})

	t.Run("empty_content", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		err := svc.Create(1, "bug", "标题", "", "")
		if err == nil {
			t.Fatal("expected error for empty content")
		}
	})

	t.Run("invalid_type_defaults_to_other", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		err := svc.Create(1, "invalid_type", "标题", "内容", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if repo.feedbacks[0].Type != "other" {
			t.Fatalf("expected type 'other', got '%s'", repo.feedbacks[0].Type)
		}
	})

	t.Run("valid_types", func(t *testing.T) {
		validTypes := []string{"bug", "suggestion", "other"}
		for _, ft := range validTypes {
			repo := newMockFeedbackRepo()
			svc := NewServiceWithDeps(repo)

			err := svc.Create(1, ft, "标题", "内容", "")
			if err != nil {
				t.Fatalf("expected no error for type '%s', got %v", ft, err)
			}
			if repo.feedbacks[0].Type != ft {
				t.Fatalf("expected type '%s', got '%s'", ft, repo.feedbacks[0].Type)
			}
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		repo.errCreate = errors.New("connection refused")
		svc := NewServiceWithDeps(repo)

		err := svc.Create(1, "bug", "标题", "内容", "")
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		feedbacks, total, err := svc.List(1, 20, nil, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total != 0 {
			t.Fatalf("expected total=0, got %d", total)
		}
		if len(feedbacks) != 0 {
			t.Fatalf("expected 0 feedbacks, got %d", len(feedbacks))
		}
	})

	t.Run("invalid_page_defaults", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		feedbacks, _, err := svc.List(0, 0, nil, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if feedbacks == nil {
			t.Fatal("expected non-nil feedbacks slice")
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		repo.errList = errors.New("timeout")
		svc := NewServiceWithDeps(repo)

		_, _, err := svc.List(1, 20, nil, "")
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestListByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		repo.feedbacks = []FeedbackResponse{
			{Feedback: Feedback{ID: 1, UserID: 1, Title: "user1"}},
			{Feedback: Feedback{ID: 2, UserID: 2, Title: "user2"}},
		}
		svc := NewServiceWithDeps(repo)

		feedbacks, total, err := svc.ListByUser(1, 1, 20)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total != 1 {
			t.Fatalf("expected total=1, got %d", total)
		}
		if len(feedbacks) != 1 {
			t.Fatalf("expected 1 feedback, got %d", len(feedbacks))
		}
	})

	t.Run("invalid_params", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		feedbacks, _, err := svc.ListByUser(1, -1, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		// nil or empty slice are both acceptable
		_ = feedbacks
	})
}

func TestUpdateStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		repo.feedbacks = []FeedbackResponse{
			{Feedback: Feedback{ID: 1, UserID: 1, Status: 0}},
		}
		svc := NewServiceWithDeps(repo)

		err := svc.UpdateStatus(1, 1, "已处理")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if repo.feedbacks[0].Status != 1 {
			t.Fatalf("expected status=1, got %d", repo.feedbacks[0].Status)
		}
	})

	t.Run("invalid_status", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		svc := NewServiceWithDeps(repo)

		err := svc.UpdateStatus(1, 5, "")
		if err == nil {
			t.Fatal("expected error for invalid status")
		}
	})

	t.Run("valid_statuses", func(t *testing.T) {
		validStatuses := []int8{0, 1, 2}
		for _, status := range validStatuses {
			repo := newMockFeedbackRepo()
			repo.feedbacks = []FeedbackResponse{
				{Feedback: Feedback{ID: 1, UserID: 1, Status: 0}},
			}
			svc := NewServiceWithDeps(repo)

			err := svc.UpdateStatus(1, status, "回复")
			if err != nil {
				t.Fatalf("expected no error for status=%d, got %v", status, err)
			}
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockFeedbackRepo()
		repo.errUpdate = errors.New("connection lost")
		svc := NewServiceWithDeps(repo)

		err := svc.UpdateStatus(1, 1, "")
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}
