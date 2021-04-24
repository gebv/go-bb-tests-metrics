test:
	go test -v -tags testrunmain -coverpkg=./... -coverprofile=c.out  ./... --run=TestCase2
	go tool cover -func=c.out
