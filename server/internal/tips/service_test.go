package tips

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// mockTipsRepo Mock 贴士仓储
type mockTipsRepo struct {
	tips    []SurvivalTip
	nextID  int64
	errGet  error
	errSave error
}

func newMockTipsRepo() *mockTipsRepo {
	return &mockTipsRepo{
		tips:   make([]SurvivalTip, 0),
		nextID: 1,
	}
}

func (m *mockTipsRepo) GetRandom() (*SurvivalTip, error) {
	if m.errGet != nil {
		return nil, m.errGet
	}
	for i := range m.tips {
		if m.tips[i].IsActive {
			return &m.tips[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockTipsRepo) GetList(page, pageSize int) ([]SurvivalTip, int64, error) {
	if m.errGet != nil {
		return nil, 0, m.errGet
	}
	var active []SurvivalTip
	for _, t := range m.tips {
		if t.IsActive {
			active = append(active, t)
		}
	}
	total := int64(len(active))
	offset := (page - 1) * pageSize
	if offset >= len(active) {
		return []SurvivalTip{}, total, nil
	}
	end := offset + pageSize
	if end > len(active) {
		end = len(active)
	}
	return active[offset:end], total, nil
}

func (m *mockTipsRepo) Create(tip *SurvivalTip) error {
	if m.errSave != nil {
		return m.errSave
	}
	tip.ID = m.nextID
	m.nextID++
	m.tips = append(m.tips, *tip)
	return nil
}

func (m *mockTipsRepo) Update(tip *SurvivalTip) error {
	if m.errSave != nil {
		return m.errSave
	}
	for i, t := range m.tips {
		if t.ID == tip.ID {
			m.tips[i] = *tip
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockTipsRepo) Delete(id int64) error {
	if m.errSave != nil {
		return m.errSave
	}
	for i, t := range m.tips {
		if t.ID == id {
			m.tips[i].IsActive = false
			return nil
		}
	}
	return nil
}

func TestGetRandom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.tips = append(repo.tips, SurvivalTip{
			ID: 1, Content: "测试贴士", Category: "mindset", IsActive: true,
		})
		svc := NewTipsServiceWithDeps(repo)

		tip, err := svc.GetRandom()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tip.Content != "测试贴士" {
			t.Fatalf("expected content '测试贴士', got '%s'", tip.Content)
		}
	})

	t.Run("empty_returns_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		svc := NewTipsServiceWithDeps(repo)

		_, err := svc.GetRandom()
		if err == nil {
			t.Fatal("expected error when no tips exist")
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errGet = errors.New("connection refused")
		svc := NewTipsServiceWithDeps(repo)

		_, err := svc.GetRandom()
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestGetList(t *testing.T) {
	t.Run("success_with_pagination", func(t *testing.T) {
		repo := newMockTipsRepo()
		for i := 1; i <= 5; i++ {
			repo.tips = append(repo.tips, SurvivalTip{
				ID: int64(i), Content: "贴士", IsActive: true,
			})
		}
		svc := NewTipsServiceWithDeps(repo)

		tips, total, err := svc.GetList(1, 3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total != 5 {
			t.Fatalf("expected total=5, got %d", total)
		}
		if len(tips) != 3 {
			t.Fatalf("expected 3 tips, got %d", len(tips))
		}
	})

	t.Run("invalid_page_defaults", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.tips = append(repo.tips, SurvivalTip{
			ID: 1, Content: "贴士", IsActive: true,
		})
		svc := NewTipsServiceWithDeps(repo)

		tips, _, err := svc.GetList(0, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tips) != 1 {
			t.Fatalf("expected 1 tip with default params, got %d", len(tips))
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errGet = errors.New("timeout")
		svc := NewTipsServiceWithDeps(repo)

		_, _, err := svc.GetList(1, 10)
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockTipsRepo()
		svc := NewTipsServiceWithDeps(repo)

		tip, err := svc.Create("新贴士", "finance")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tip.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if tip.Content != "新贴士" {
			t.Fatalf("expected content '新贴士', got '%s'", tip.Content)
		}
		if !tip.IsActive {
			t.Fatal("expected tip to be active")
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errSave = errors.New("disk full")
		svc := NewTipsServiceWithDeps(repo)

		_, err := svc.Create("贴士", "life")
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.tips = append(repo.tips, SurvivalTip{
			ID: 1, Content: "旧内容", Category: "life", IsActive: true,
		})
		svc := NewTipsServiceWithDeps(repo)

		err := svc.Update(1, "新内容", "mindset")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errSave = errors.New("connection lost")
		svc := NewTipsServiceWithDeps(repo)

		err := svc.Update(1, "内容", "life")
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.tips = append(repo.tips, SurvivalTip{
			ID: 1, Content: "贴士", IsActive: true,
		})
		svc := NewTipsServiceWithDeps(repo)

		err := svc.Delete(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errSave = errors.New("permission denied")
		svc := NewTipsServiceWithDeps(repo)

		err := svc.Delete(1)
		if err == nil {
			t.Fatal("expected error on db failure")
		}
	})
}
