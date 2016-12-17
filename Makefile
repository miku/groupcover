groupcover: cmd/groupcover/main.go rewriter.go
	goimports -w .
	go build -o groupcover cmd/groupcover/main.go

clean:
	rm -f groupcover
