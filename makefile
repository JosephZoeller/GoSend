all: release_linux

build_linux:
	go build ./cmd/typetest
	go build ./cmd/typetestd

release_linux:
	mkdir -p ./Release_Build/linux-amd64/server/web
	go build -o ./Release_Build/linux-amd64/GMG ./client
	go build -o ./Release_Build/linux-amd64/server/GMG-Server ./server
	cp ./web/tables.html ./Release_Build/linux-amd64/server/web/tables.html

release_windows:
	mkdir -p ./Release_Build/windows-amd64/server/web/
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/GMG.exe ./client
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/server/GMG-Server.exe ./server
	cp ./web/tables.html ./Release_Build/windows-amd64/server/web/tables.html

release: release_linux release_windows