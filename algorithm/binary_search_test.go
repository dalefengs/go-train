package algorithm

import (
	"fmt"
	"testing"
)

func binarySearch(arr []int64, target int64) (int64, int) {
	arrLen := len(arr)
	if arrLen == 0 {
		return -1, 0
	}
	left := 0
	right := arrLen - 1
	count := 0
	for left <= right {
		count++
		middle := left + (right-left)/2
		if arr[middle] == target {
			return int64(middle), count
		}
		if arr[middle] < target {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}
	return -1, count
}

func TestBinarySearch(t *testing.T) {
	testCases := []struct {
		arr    []int64
		target int64
		index  int64
		count  int // 预期的查找次数
	}{
		// 测试用例1：目标值存在于数组中
		{
			arr:    []int64{1, 3, 5, 7, 9},
			target: 5,
			index:  2,
			count:  1,
		},
		// 无数据
		{
			arr:    []int64{1, 3, 5, 7, 9},
			target: 4,
			count:  3,
			index:  -1,
		},
		// 测试用例3：目标值在数组的第一个元素
		{
			arr:    []int64{1, 3, 5, 7, 9},
			target: 1,
			count:  2,
			index:  0,
		},
		// 测试用例4：目标值在数组的最后一个元素
		{
			arr:    []int64{1, 3, 5, 7, 9},
			target: 9,
			count:  3,
			index:  4,
		},
		// 测试用例5：空数组
		{
			arr:    []int64{},
			target: 5,
			count:  0,
			index:  -1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("target=%d", tc.target), func(t *testing.T) {
			index, count := binarySearch(tc.arr, tc.target)
			if index != tc.index {
				t.Errorf("expected %v, but got %v", tc.index, index)
			}
			fmt.Println("遍历次数为", count)
			//if count != tc.count {
			//	t.Errorf("expected count %d, but got %d", tc.count, count)
			//}
		})
	}
}
