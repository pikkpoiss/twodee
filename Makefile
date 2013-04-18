.PHONY: cprofile-basic mprofile-basic wprofile-basic

example-basic: examples/basic/*.go
	go build -o=$@ examples/basic/*.go

example-basic.cprofile: example-basic
	./example-basic -cpuprofile=$@

example-basic.mprofile: example-basic
	./example-basic -memprofile=$@

cprofile-basic: example-basic example-basic.cprofile
	go tool pprof --web example-basic example-basic.cprofile

mprofile-basic: example-basic example-basic.mprofile
	go tool pprof --web example-basic example-basic.mprofile

wprofile-basic: example-basic
	./example-basic -webprofile &
	echo "http://localhost:6060/debug/pprof/"
	echo "go tool pprof http://localhost:6060/debug/pprof/heap"
