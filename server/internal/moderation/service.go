package moderation

import (
	"strings"
	"unicode"
)

// ModerationService 内容审核服务
type ModerationService struct {
	// 内置敏感词库（MVP 阶段使用内置词库，后续可从数据库加载）
	builtinWords []string
}

// NewModerationService 创建审核服务
func NewModerationService() *ModerationService {
	return &ModerationService{
		builtinWords: defaultSensitiveWords(),
	}
}

// CheckContent 检查内容是否包含敏感词
func (s *ModerationService) CheckContent(text string) *CheckResult {
	if text == "" {
		return &CheckResult{Passed: true}
	}

	// 预处理：转小写，去除特殊字符
	normalized := normalize(text)

	var hitWords []string
	for _, word := range s.builtinWords {
		if contains(normalized, word) {
			hitWords = append(hitWords, word)
		}
	}

	if len(hitWords) > 0 {
		return &CheckResult{
			Passed:  false,
			Words:   hitWords,
			Message: "内容包含违规词汇，请修改后重试",
		}
	}

	return &CheckResult{Passed: true}
}

// CheckTitle 检查标题
func (s *ModerationService) CheckTitle(title string) *CheckResult {
	return s.CheckContent(title)
}

// AddWord 添加敏感词到运行时词库（可用于管理后台动态添加）
func (s *ModerationService) AddWord(word string) {
	s.builtinWords = append(s.builtinWords, normalize(word))
}

// normalize 文本预处理
func normalize(text string) string {
	// 转小写
	text = strings.ToLower(text)
	// 去除空白字符
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)
	return text
}

// contains 检查文本是否包含关键词（支持变体绕过检测）
func contains(text, word string) bool {
	// 直接匹配
	if strings.Contains(text, word) {
		return true
	}

	// 去除常见分隔符后匹配
	separated := strings.NewReplacer("*", "", " ", "", ".", "", "-", "", "_", "")
	cleaned := separated.Replace(text)
	if strings.Contains(cleaned, word) {
		return true
	}

	return false
}

// defaultSensitiveWords 返回默认敏感词库（MVP 阶段使用基础词库）
func defaultSensitiveWords() []string {
	return []string{
		// 政治敏感
		"习近平", "毛泽东", "六四", "天安门事件",
		// 色情
		"色情", "裸体", "性交",
		// 暴力
		"杀人", "炸弹制作", "恐怖袭击",
		// 赌博
		"赌博", "赌场", "博彩",
		// 毒品
		"毒品", "冰毒", "大麻",
		// 诈骗
		"诈骗", "传销", "非法集资",
		// 歧视
		"黑鬼", "支那",
		// 其他违法
		"假证", "办证", "代孕",
	}
}
