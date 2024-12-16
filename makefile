run/agent: 
	go run cmd/agent/main.go

run/server: 
	go run cmd/server/main.go

build/server:
	cd cmd/server && \
		go build -buildvcs=false -o server && \
		cd ../..

build/agent:
	cd cmd/agent && \
		go build -buildvcs=false -o agent && \
		cd ../..

run/test1: build/server
	metricstest -test.v -test.run="^TestIteration1$$" \
		-binary-path=cmd/server/server

run/test2: build/agent
	metricstest -test.v -test.run="^TestIteration2[AB]*$$" \
		-source-path=. \
 		-agent-binary-path=cmd/agent/agent

run/tests:
	go test ./... -count=1 -coverprofile ./profiles/cover.out && go tool cover -func ./profiles/cover.out

show/cover:
	go tool cover -html=./profiles/cover.out

PHONY: run/test1, run/test2, build/server, run/tests, show/cover, run/agent