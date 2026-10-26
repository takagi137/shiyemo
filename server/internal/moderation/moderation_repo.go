package moderation

// ModerationRepository 审核仓储接口
type ModerationRepository interface {
	CreateWord(word *SensitiveWord) error
	GetWordByID(id int64) (*SensitiveWord, error)
	GetWordByName(name string) (*SensitiveWord, error)
	ListWords(page, size int) ([]SensitiveWord, int64, error)
	DeleteWord(id int64) error
	GetAllWords() ([]string, error)
}
