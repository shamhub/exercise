PHONY: run build clean test

run:
	go run cmd/weather/main.go


build:
	go build  -o main cmd/weather/main.go

test:
	go test -v -race -tags=unit -cover ./extract

clean:
	rm main