package util

import (
	"errors"
)

func Get_max_index(sl []int) (int, int) {
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

func Get_sorted_indexes(sl []int) []int {
	slcopy := make([]int, len(sl))
	copy(slcopy, sl)
	min_int := -int(^uint(0)>>1) - 1
	res := []int{}
	length := len(slcopy)
	for i := 0; i < length; i++ {
		_, idx := Get_max_index(slcopy)
		res = append(res, idx)
		slcopy[idx] = min_int
	}
	return res
}

func Get_element_sum(sl []int, inds []int) []int {
	sum := 0
	res := []int{}
	for _, ind := range inds {
		sum += sl[ind]
		res = append(res, sum)
	}
	return res
}

func Necessary_inputs(utxos []map[string]string, amount int) ([]map[string]string, error) {
	res := make([]map[string]string, 0)
	values := []int{}
	for _, uo := range utxos {
		values = append(values, Str_to_int(uo["value"]))
	}
	values_index_desc := Get_sorted_indexes(values)
	necessary_indexes := []int{}
	total := 0
	for i := 0; total < amount && i < len(values_index_desc); i++ {
		total += values[values_index_desc[i]]
		necessary_indexes = append(necessary_indexes, values_index_desc[i])
	}
	if total < amount {
		return nil, errors.New("残高が不足しています。")
	}
	for _, i := range necessary_indexes {
		res = append(res, utxos[i])
	}
	return res, nil
}
