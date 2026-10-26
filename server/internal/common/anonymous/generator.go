package anonymous

import (
	"fmt"
	"math/rand"
	"time"
)

// 匿名名称组成部分
var (
	adjectives = []string{
		"快乐的", "忧伤的", "迷茫的", "坚强的", "温柔的",
		"勇敢的", "沉默的", "乐观的", "焦虑的", "平静的",
		"孤独的", "自由的", "疲惫的", "充满希望的", "失落的",
	}

	animals = []string{
		"小猫", "小狗", "小兔", "小熊", "小鹿",
		"小鸟", "小鱼", "小松鼠", "小狐狸", "小熊猫",
		"小海豚", "小企鹅", "小刺猬", "小仓鼠", "小乌龟",
	}
)

// Generator 匿名名称生成器
type Generator struct {
	rng *rand.Rand
}

// NewGenerator 创建生成器
func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate 生成匿名名称
func (g *Generator) Generate() string {
	adj := adjectives[g.rng.Intn(len(adjectives))]
	animal := animals[g.rng.Intn(len(animals))]
	return fmt.Sprintf("%s%s", adj, animal)
}

// 全局生成器实例
var defaultGenerator = NewGenerator()

// Generate 使用全局生成器生成匿名名称
func Generate() string {
	return defaultGenerator.Generate()
}
