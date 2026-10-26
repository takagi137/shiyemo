package moderation

// SensitiveWord 敏感词
type SensitiveWord struct {
	ID   int64  `json:"id"`
	Word string `json:"word"`
}

// TableName 表名
func (SensitiveWord) TableName() string {
	return "sensitive_words"
}

// CheckResult 检查结果
type CheckResult struct {
	Passed  bool     `json:"passed"`   // 是否通过
	Words   []string `json:"words"`    // 命中的敏感词列表
	Message string   `json:"message"`  // 提示信息
}
