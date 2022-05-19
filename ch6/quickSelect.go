package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSelect(nums []int, k int) int {
	rand.Seed(time.Now().UnixNano())
	return quickSort(nums, k, 0, len(nums)-1)
}

func quickSort(nums []int, k int, low int, high int) int {
	mid := divide(nums, low, high)
	if mid == k-1 {
		return nums[mid]
	} else if mid < k-1 {
		return quickSort(nums, k, mid+1, high)
	} else {
		return quickSort(nums, k, low, mid-1)
	}
}

func divide(nums []int, low int, high int) int {
	i := rand.Intn(high-low+1) + low
	nums[low], nums[i] = nums[i], nums[low]
	pixos := nums[low]
	for low < high {
		for low < high && nums[high] > pixos {
			high--
		}
		if low < high {
			nums[low] = nums[high]
			low++
		}
		for low < high && nums[low] < pixos {
			low++
		}
		if low < high {
			nums[high] = nums[low]
			high--
		}
	}
	nums[low] = pixos
	return low
}

func main() {
	nums := []int{9, 199, 32, 48, 13, -1, 43, 165, -32}
	fmt.Println(quickSelect(nums, 3))
	fmt.Println(nums)
}
