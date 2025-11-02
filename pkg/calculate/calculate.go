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
	for targetAmount := orderedItems; targetAmount < orderedItems+largestPackSize; targetAmount++ {
		solution := findOptimalCombination(targetAmount, packSizes)
		if solution != nil {
			return solution, nil
		}
	}

	return nil, ErrCannotFulfill
}

func findOptimalCombination(targetAmount int64, packSizes []int64) map[int64]int64 {
	minPacksNeeded := make([]int64, targetAmount+1)

	lastPackUsed := make([]int64, targetAmount+1)

	for amount := 1; int64(amount) <= targetAmount; amount++ {
		minPacksNeeded[amount] = targetAmount + 1
		lastPackUsed[amount] = -1
	}

	for currentAmount := 1; int64(currentAmount) <= targetAmount; currentAmount++ {
		for _, packSize := range packSizes {
			remainingAmount := int64(currentAmount) - packSize
			if remainingAmount >= 0 && minPacksNeeded[remainingAmount]+1 < minPacksNeeded[currentAmount] {
				minPacksNeeded[currentAmount] = minPacksNeeded[remainingAmount] + 1
				lastPackUsed[currentAmount] = packSize
			}
		}
	}

	if minPacksNeeded[targetAmount] > targetAmount {
		return nil
	}

	packCounts := make(map[int64]int64)
	remainingAmount := targetAmount

	for remainingAmount > 0 {
		packSize := lastPackUsed[remainingAmount]
		packCounts[packSize]++
		remainingAmount -= packSize
	}

	return packCounts
}
