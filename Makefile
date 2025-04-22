.PHONY: test
test:
	@go test

.PHONY: examples
examples:
	@go build -o examples/bin examples/layers_accesskey.go 
	@go build -o examples/bin examples/layers_apikey.go
	@go build -o examples/bin examples/layout_accesskey.go
	@go build -o examples/bin examples/meeting.go
	@go build -o examples/bin examples/observer.go
	@go build -o examples/bin examples/shutdown.go
	@go build -o examples/bin examples/snapshots.go
	@go build -o examples/bin examples/webhook-listener.go
