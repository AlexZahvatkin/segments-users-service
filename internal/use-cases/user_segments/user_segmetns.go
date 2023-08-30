package usecases_user_segments

import (
	"errors"
	"math/rand"
)

func PickRandomIds(percent float64, ids []int64) ([]int64, error) {
	if percent > 100 {
		return nil, errors.New("Wrong percent value")
	}
	var res []int64
	for _, item := range ids {
		if rand.Float64() * 100 < percent {
			res = append(res, item)
		}
	}
	return res, nil
}