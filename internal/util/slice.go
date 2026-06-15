package util

import (
	"errors"
)

func GetMaxIndex(sl []int) (int, int) {
	max := sl[0]
	ci := 0
	for i, num := range sl {
		if num > max {
			max = num
			ci = i
		}
	}
	return max, ci
}

func GetSortedIndexes(sl []int) []int {
	slcopy := make([]int, len(sl))
	copy(slcopy, sl)
	minInt := -int(^uint(0)>>1) - 1
	res := []int{}
	length := len(slcopy)
	for i := 0; i < length; i++ {
		_, idx := GetMaxIndex(slcopy)
		res = append(res, idx)
		slcopy[idx] = minInt
	}
	return res
}

func GetElementSum(sl []int, inds []int) []int {
	sum := 0
	res := []int{}
	for _, ind := range inds {
		sum += sl[ind]
		res = append(res, sum)
	}
	return res
}

func NecessaryInputs(utxos []map[string]string, amount int) ([]map[string]string, error) {
	res := make([]map[string]string, 0)
	values := []int{}
	for _, uo := range utxos {
		values = append(values, StrToInt(uo["value"]))
	}
	valuesIndexDesc := GetSortedIndexes(values)
	necessaryIndexes := []int{}
	total := 0
	for i := 0; total < amount && i < len(valuesIndexDesc); i++ {
		total += values[valuesIndexDesc[i]]
		necessaryIndexes = append(necessaryIndexes, valuesIndexDesc[i])
	}
	if total < amount {
		return nil, errors.New("残高が不足しています。")
	}
	for _, i := range necessaryIndexes {
		res = append(res, utxos[i])
	}
	return res, nil
}
