build/server:
	cd cmd/server && \
		go build -buildvcs=false -o server && \
		cd ../..

run/test1: build/server
	metricstest -test.v -test.run="^TestIteration1$$" \
		-binary-path=cmd/server/server

PHONY: run/test1, build/server