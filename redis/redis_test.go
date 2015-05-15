package redis

import (
	"testing"
)

const (
	key   string = "redis-bloom-filter-test"
	value string = "data"
)

func Test_Badd(t *testing.T) {

	Badd(key, value)

}

func Test_Bexist(t *testing.T) {
	exist := Bexist(key, value)
	if exist {
		t.Log("bloom Badd success also success on Bexist")
	} else {
		t.Error("bloom method error")
	}
}
