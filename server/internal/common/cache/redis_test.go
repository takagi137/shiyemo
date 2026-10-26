package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleflight_Basic(t *testing.T) {
	sf := NewSingleflight()

	result, err := sf.Do("key1", func() (interface{}, error) {
		return "value1", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "value1" {
		t.Fatalf("expected value1, got %v", result)
	}
}

func TestSingleflight_Error(t *testing.T) {
	sf := NewSingleflight()

	expectedErr := errors.New("db error")
	_, err := sf.Do("key2", func() (interface{}, error) {
		return nil, expectedErr
	})
	if err != expectedErr {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestSingleflight_Dedup(t *testing.T) {
	sf := NewSingleflight()

	var callCount int32
	var wg sync.WaitGroup
	concurrency := 10

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			sf.Do("dedup-key", func() (interface{}, error) {
				atomic.AddInt32(&callCount, 1)
				time.Sleep(50 * time.Millisecond)
				return "result", nil
			})
		}()
	}
	wg.Wait()

	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("expected fn called 1 time, got %d", callCount)
	}
}

func TestSingleflight_DifferentKeys(t *testing.T) {
	sf := NewSingleflight()

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	wg.Add(2)
	go func() {
		defer wg.Done()
		v, _ := sf.Do("keyA", func() (interface{}, error) {
			time.Sleep(30 * time.Millisecond)
			return "A", nil
		})
		mu.Lock()
		results["A"] = v.(string)
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		v, _ := sf.Do("keyB", func() (interface{}, error) {
			time.Sleep(30 * time.Millisecond)
			return "B", nil
		})
		mu.Lock()
		results["B"] = v.(string)
		mu.Unlock()
	}()
	wg.Wait()

	if results["A"] != "A" || results["B"] != "B" {
		t.Fatalf("expected A and B, got %v", results)
	}
}

func TestSingleflight_SequentialReuse(t *testing.T) {
	sf := NewSingleflight()

	// 第一次调用
	v1, _ := sf.Do("reuse", func() (interface{}, error) {
		return 1, nil
	})
	// 第二次调用（key 已清除，应重新执行 fn）
	v2, _ := sf.Do("reuse", func() (interface{}, error) {
		return 2, nil
	})

	if v1.(int) != 1 || v2.(int) != 2 {
		t.Fatalf("expected 1 then 2, got %v then %v", v1, v2)
	}
}
