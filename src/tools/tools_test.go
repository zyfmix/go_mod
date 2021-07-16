package tools

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestTools(t *testing.T) {
	numbers := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6}
	fmt.Println("Numbers", zap.Any("numbers", numbers))
	fmt.Println("Numbers", zap.Any("contain 5", ContainInt(numbers, 5)))
	fmt.Println("Numbers", zap.Any("contain 8", ContainInt(numbers, 8)))
	fmt.Println("Numbers", zap.Any("deduped", Dedup(numbers)))

	// 求最小最大值...
	fmt.Println("(3, 9, 6, 2)", zap.Any("MaxOf", MinOf(3, 9, 6, 2)))
	fmt.Println("(3, 9, 6, 2)", zap.Any("MinOf", MinOf(3, 9, 6, 2)))
}
