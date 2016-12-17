groupcover: cmd/groupcover/main.go rewriter.go
	go build -o groupcover cmd/groupcover/main.go

clean:
	rm -f groupcover
