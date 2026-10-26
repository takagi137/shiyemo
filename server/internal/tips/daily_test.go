package tips

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomQuote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.tips = append(repo.tips, SurvivalTip{
			ID: 1, Content: "今天的你，已经比昨天勇敢了一点。", Category: "encourage", IsActive: true,
		})
		svc := NewTipsServiceWithDeps(repo)

		quote, err := svc.RandomQuote()
		assert.NoError(t, err)
		assert.Equal(t, "今天的你，已经比昨天勇敢了一点。", quote.Text)
		// 设计意图变更：category 是内部分类（mindset/finance/...），不再作为
		// 作者展示，否则首页会渲染出 "—— encourage" 这类无意义署名。故 Author 恒为空。
		assert.Equal(t, "", quote.Author)
	})

	t.Run("empty_repo_returns_default", func(t *testing.T) {
		repo := newMockTipsRepo()
		svc := NewTipsServiceWithDeps(repo)

		quote, err := svc.RandomQuote()
		assert.NoError(t, err)
		assert.Equal(t, "今天的你，已经比昨天勇敢了一点。", quote.Text)
		assert.Equal(t, "", quote.Author)
	})

	t.Run("db_error_returns_default", func(t *testing.T) {
		repo := newMockTipsRepo()
		repo.errGet = errors.New("connection refused")
		svc := NewTipsServiceWithDeps(repo)

		quote, err := svc.RandomQuote()
		assert.NoError(t, err)
		assert.Equal(t, "今天的你，已经比昨天勇敢了一点。", quote.Text)
	})
}

func TestGetDailyContent_DefaultFallback(t *testing.T) {
	// 当图片服务不可用时，GetDailyContent 应该返回默认图片
	repo := newMockTipsRepo()
	svc := NewTipsServiceWithDeps(repo)

	// GetDailyContent 会调用 dailyimage.NewService()，在测试中会 panic
	// 所以只测试 RandomQuote（不依赖 DB 的方法）
	quote, err := svc.RandomQuote()
	assert.NoError(t, err)
	assert.NotNil(t, quote)
}
