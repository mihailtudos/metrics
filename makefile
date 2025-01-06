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

run/test3:
	metricstest -test.v -test.run=^TestIteration3[AB]*$ \
    	-source-path=. \
    	-agent-binary-path=cmd/agent/agent \
    	-binary-path=cmd/server/server
	
run/test4:
	ADDRESS=localhost:8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration4$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent

run/test5:
	ADDRESS=localhost:8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration5$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent

run/test6:
	ADDRESS=localhost:8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration6$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent \

run/test7:
	ADDRESS=localhost:8080 TEMP_FILE=out.txt metricstest -test.v -test.run=^TestIteration7$$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=. \
        -agent-binary-path=cmd/agent/agent \

run/test8:
	TEMP_FILE=out.txt metricstest -test.v -test.run=^TestIteration8$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=.-agent-binary-path=cmd/agent/agent \

autotest/run9:
	TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration9$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
        -server-port=8080 \
        -source-path=. \
        -file-storage-path=/tmp/metrics-db.json \
        -agent-binary-path=cmd/agent/agent \

autotest/run10: db/run
	 SERVER_PORT=8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration10[AB]$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://metrics:metrics@localhost:5432/metrics?sslmode=disable' \
        -server-port=8080 \
        -source-path=.

autotest/run11: db/run
	 SERVER_PORT=8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration11$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://metrics:metrics@localhost:5432/metrics?sslmode=disable' \
        -server-port=8080 \
        -source-path=.

autotest/run12: db/run
	 SERVER_PORT=8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration12$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://metrics:metrics@localhost:5432/metrics?sslmode=disable' \
        -server-port=8080 \
        -source-path=.

autotest/run13: db/run
	 SERVER_PORT=8080 TEMP_FILE=out.txt metricstest -test.v -test.run="^TestIteration13$$" \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://metrics:metrics@localhost:5432/metrics?sslmode=disable' \
        -server-port=8080 \
        -source-path=.
		
run/tests:
	go test ./... -count=1 -coverprofile ./profiles/cover.out && go tool cover -func ./profiles/cover.out

show/cover:
	go tool cover -html=./profiles/cover.out

run/vet:
	go vet -vettool=$(which statictest) ./...

PHONY: run/test1, run/test2, run/test3, run/test4, run/test5, run/test6, run/test7, run/test8, build/server, run/tests, show/cover, run/agent, run/vet