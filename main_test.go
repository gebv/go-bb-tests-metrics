package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCase2(t *testing.T) {
	handler("a", "a")
	handler("a", "b")
	handler("a", "c")
	handler("b", "a")
	handler("b", "b")
	handler("b", "c")
}

func TestCase1(t *testing.T) {
	go main()
	time.Sleep(time.Millisecond * 100)

	res, err := http.Get("http://127.0.0.1:8080/differentCase")
	assert.NoError(t, err)
	assertBody(t, res, "case2")

	res, err = http.Get("http://127.0.0.1:8080/differentCase?foo=bar")
	assert.NoError(t, err)
	assertBody(t, res, "foobar")

	res, err = http.Get("http://127.0.0.1:8080/onlyPost")
	assert.NoError(t, err)
	assertBody(t, res, "only POST")

	res, err = http.Post("http://127.0.0.1:8080/onlyPost", "application/json", bytes.NewReader([]byte(`{}`)))
	assert.NoError(t, err)
	assertBody(t, res, "hello")
}

func assertBody(t *testing.T, res *http.Response, want string) {
	t.Helper()
	dat, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	err = res.Body.Close()
	assert.NoError(t, err)
	assert.EqualValues(t, want, string(dat))
}
