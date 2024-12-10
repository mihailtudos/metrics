build/server:
	cd cmd/server && \
		go build -buildvcs=false -o server && \
		cd ../..

run/test1: build/server
	metricstest -test.v -test.run="^TestIteration1$$" \
		-binary-path=cmd/server/server

run/tests:
	go test ./... -count=1 -coverprofile ./profiles/cover.out && go tool cover -func ./profiles/cover.out

show/cover:
	go tool cover -html=./profiles/cover.out

PHONY: run/test1, build/server, run/tests, show/cover