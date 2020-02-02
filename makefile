all: release_linux

build_linux:
	go build ./cmd/client
	go build ./cmd/rvprox
	go build ./cmd/server

release_linux:
	mkdir -p ./Release_Build/linux-amd64/server/web
	mkdir -p ./Release_Build/linux-amd64/rvprox
	go build -o ./Release_Build/linux-amd64/GMG ./cmd/client
	go build -o ./Release_Build/linux-amd64/reverseProxy/GMG-ReverseProxy ./cmd/rvprox
	go build -o ./Release_Build/linux-amd64/server/GMG-Server ./cmd/server
	cp ./web/tables.html ./Release_Build/linux-amd64/server/web/tables.html

release_windows:
	mkdir -p ./Release_Build/windows-amd64/server/web/
	mkdir -p ./Release_Build/windows-amd64/rvprox
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/GMG.exe ./cmd/client
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/reverseProxy/GMG-ReverseProxy.exe ./cmd/rvprox
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/server/GMG-Server.exe ./cmd/server
	cp ./web/tables.html ./Release_Build/windows-amd64/server/web/tables.html

release: release_linux release_windows