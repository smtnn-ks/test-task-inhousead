run: 
	go run ./cmd

build: 
	go build -o ./bin/app ./cmd

build_debug: 
	go build -gcflags 'all=-N -l' -o ./bin/app-debug ./cmd

debug: build_debug
	gdlv exec ./bin/app-debug
