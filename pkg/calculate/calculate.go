package calculate

import (
	"errors"
	"slices"
)

var (
	ErrEmptySizes    = errors.New("sizes can't be empty")
	ErrCannotFulfill = errors.New("can't fulfill the order with current packs")
)

// CalculatePackages returns the number of packages of each size to fulfill the order.
func CalculatePackages(orderedItems int64, packSizes []int64) (map[int64]int64, error) {
	if len(packSizes) == 0 {
		return nil, ErrEmptySizes
	}

	largestPackSize := slices.Max(packSizes)
	solution := findOptimalCombination(orderedItems, largestPackSize, packSizes)
	if solution != nil {
		return solution, nil
	}

	return nil, ErrCannotFulfill
}

type state struct {
	amount    int64
	packCount int64
}

func findOptimalCombination(orderedItems int64, largestPackSize int64, packSizes []int64) map[int64]int64 {
	maxAmount := orderedItems + largestPackSize

	type parentInfo struct {
		previousAmount int64
		packUsed       int64
	}
	parent := make(map[int64]parentInfo)
	packCount := make(map[int64]int64)

	queue := []state{{amount: 0, packCount: 0}}
	packCount[0] = 0

	var bestAmount int64 = -1
	var bestPackCount int64 = maxAmount

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, packSize := range packSizes {
			nextAmount := current.amount + packSize
			nextPackCount := current.packCount + 1

			if nextAmount > maxAmount {
				continue
			}

			if prevCount, exists := packCount[nextAmount]; exists && prevCount <= nextPackCount {
				continue
			}

			packCount[nextAmount] = nextPackCount
			parent[nextAmount] = parentInfo{
				previousAmount: current.amount,
				packUsed:       packSize,
			}

			if nextAmount >= orderedItems {
				overshoot := nextAmount - orderedItems

				if bestAmount == -1 ||
					overshoot < (bestAmount-orderedItems) ||
					(overshoot == (bestAmount-orderedItems) && nextPackCount < bestPackCount) {
					bestAmount = nextAmount
					bestPackCount = nextPackCount
				}
			} else {
				queue = append(queue, state{amount: nextAmount, packCount: nextPackCount})
			}
		}
	}

	if bestAmount == -1 {
		return nil
	}

	result := make(map[int64]int64)
	currentAmount := bestAmount

	for currentAmount > 0 {
		info := parent[currentAmount]
		result[info.packUsed]++
		currentAmount = info.previousAmount
	}

	return result
}
