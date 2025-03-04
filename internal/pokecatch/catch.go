package pokecatch

import (
	"math/rand"
)

func SuccessfulCatch(baseExperience int) bool {
	chances := rand.Intn(baseExperience * 2)
	if chances >= baseExperience {
		return true
	}
	return false
}
