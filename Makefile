build:
	go build ./bus ./client ./rhino ./lion

test: build
	go test ./bus/...  ./client/... ./lion/... ./rhino/...
