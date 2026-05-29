.PHONY: run
run:
	@echo "🚀 Running..."
	nodemon --exec go run main.go --signal SIGTERM