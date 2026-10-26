package moderation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModerationService_CheckContent(t *testing.T) {
	service := NewModerationService()

	t.Run("empty_content", func(t *testing.T) {
		result := service.CheckContent("")
		assert.True(t, result.Passed)
	})

	t.Run("safe_content", func(t *testing.T) {
		result := service.CheckContent("今天天气真好，适合出去走走")
		assert.True(t, result.Passed)
		assert.Empty(t, result.Words)
	})

	t.Run("contains_sensitive_word", func(t *testing.T) {
		result := service.CheckContent("这个赌博网站很赚钱")
		assert.False(t, result.Passed)
		assert.Contains(t, result.Words, "赌博")
		assert.NotEmpty(t, result.Message)
	})

	t.Run("multiple_sensitive_words", func(t *testing.T) {
		result := service.CheckContent("赌博和毒品都是违法的")
		assert.False(t, result.Passed)
		assert.Len(t, result.Words, 2)
		assert.Contains(t, result.Words, "赌博")
		assert.Contains(t, result.Words, "毒品")
	})

	t.Run("case_insensitive", func(t *testing.T) {
		result := service.CheckContent("这是一个测试色情内容")
		assert.False(t, result.Passed)
		assert.Contains(t, result.Words, "色情")
	})

	t.Run("with_separators", func(t *testing.T) {
		// 测试带分隔符的变体
		result := service.CheckContent("这个赌*博网站")
		assert.False(t, result.Passed)
	})

	t.Run("safe_job_seeking_content", func(t *testing.T) {
		// 正常求职内容不应被拦截
		content := "我是一名前端工程师，3年经验，熟悉 React 和 Vue，正在寻找新的工作机会"
		result := service.CheckContent(content)
		assert.True(t, result.Passed)
	})

	t.Run("safe_emotion_content", func(t *testing.T) {
		// 正常情感倾诉不应被拦截
		content := "今天被裁员了，心情很低落，但是我相信明天会更好"
		result := service.CheckContent(content)
		assert.True(t, result.Passed)
	})
}

func TestModerationService_CheckTitle(t *testing.T) {
	service := NewModerationService()

	t.Run("safe_title", func(t *testing.T) {
		result := service.CheckTitle("求职：前端工程师")
		assert.True(t, result.Passed)
	})

	t.Run("sensitive_title", func(t *testing.T) {
		result := service.CheckTitle("赌博技巧分享")
		assert.False(t, result.Passed)
	})
}

func TestModerationService_AddWord(t *testing.T) {
	service := NewModerationService()

	// 添加自定义敏感词
	service.AddWord("自定义敏感词")

	result := service.CheckContent("这里包含自定义敏感词")
	assert.False(t, result.Passed)
	assert.Contains(t, result.Words, "自定义敏感词")
}

func TestNormalize(t *testing.T) {
	t.Run("lowercase", func(t *testing.T) {
		result := normalize("Hello World")
		assert.Equal(t, "helloworld", result)
	})

	t.Run("remove_spaces", func(t *testing.T) {
		result := normalize("你 好 世 界")
		assert.Equal(t, "你好世界", result)
	})
}

func TestContains(t *testing.T) {
	t.Run("direct_match", func(t *testing.T) {
		assert.True(t, contains("赌博网站", "赌博"))
	})

	t.Run("no_match", func(t *testing.T) {
		assert.False(t, contains("正常内容", "赌博"))
	})

	t.Run("with_separator", func(t *testing.T) {
		assert.True(t, contains("赌*博", "赌博"))
	})
}
