PHONY: run build clean test

runuser:
	go run -race cmd/user/main.go


build:
	go build  -o main cmd/user/main.go

test:
	go test -v -race -tags=unit -cover ./compute

clean:
	rm main