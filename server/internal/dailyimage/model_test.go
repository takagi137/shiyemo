package dailyimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDailyImageModel(t *testing.T) {
	img := DailyImage{
		ID:          1,
		URL:         "https://images.unsplash.com/test.jpg",
		Description: "测试图片",
		Credit:      "Test Credit",
		SortOrder:   1,
		IsActive:    true,
	}

	assert.Equal(t, int64(1), img.ID)
	assert.Equal(t, "https://images.unsplash.com/test.jpg", img.URL)
	assert.Equal(t, "测试图片", img.Description)
	assert.Equal(t, 1, img.SortOrder)
	assert.True(t, img.IsActive)
}
