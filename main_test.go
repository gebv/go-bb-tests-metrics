package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func TestCase2(t *testing.T) {
// 	handler("a", "a")
// 	handler("a", "b")
// 	handler("a", "c")
// 	handler("b", "a")
// 	handler("b", "b")
// 	handler("b", "c")
// }

var numIter = flag.Int64("num-iter", 0, "")

func TestMain(m *testing.T) {
	flag.Parse()
}

func TestCases(t *testing.T) {
	go main()
	time.Sleep(time.Millisecond * 100)
	cases := [][]string{
		{"a", "a", "a"},
		{"a", "b", "b"},
		{"a", "c", "c - not matched"},
		{"b", "a", "not exists"},
		{"b", "b", "not exists"},
		{"b", "c", "not exists"},
	}
	for _, case_ := range cases[0:*numIter] {
		do(case_[0], case_[1], case_[2])(t)
	}
}

func do(a, b, want string) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := http.Get("http://127.0.0.1:8080/handler?a=" + a + "&b=" + b)
		assert.NoError(t, err)
		assertBody(t, res, want)
	}
}

func assertBody(t *testing.T, res *http.Response, want string) {
	t.Helper()
	dat, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	err = res.Body.Close()
	assert.NoError(t, err)
	assert.EqualValues(t, want, string(dat))
}
