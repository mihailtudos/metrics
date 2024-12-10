run/test1:
	metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/server/server

PHONY: run/test1