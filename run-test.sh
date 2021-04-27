#!/bin/bash

go test -v -tags testrunmain -coverpkg=./... -coverprofile=c.out  ./... --run=TestCases --num-iter $1
go tool cover -func=c.out
go tool cover -html=c.out -o coverage-test$1.html
