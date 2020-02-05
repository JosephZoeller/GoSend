all: test_linux

test_linux:
	go build -o ./test/client ./cmd/client
	go build -o ./test/rvprox ./cmd/rvprox
	go build -o ./test/server ./cmd/server
	go build -o ./test/logmgr ./cmd/logmgr
	cp ./scripts/networktest.sh ./test/t.sh

release_linux:
	go build -o ./Release_Build/linux-amd64/client ./cmd/GMG-Client.exe
	go build -o ./Release_Build/linux-amd64/reverseProxy/GMG-ReverseProxy ./cmd/rvprox
	go build -o ./Release_Build/linux-amd64/server/GMG-Server ./cmd/server
	#mkdir -p ./Release_Build/linux-amd64/server/web
	#cp ./web/tables.html ./Release_Build/linux-amd64/server/web/tables.html

release_windows:
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/GMG-Client.exe ./cmd/client
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/reverseProxy/GMG-ReverseProxy.exe ./cmd/rvprox
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/server/GMG-Server.exe ./cmd/server
	#mkdir -p ./Release_Build/windows-amd64/server/web/
	#cp ./web/tables.html ./Release_Build/windows-amd64/server/web/tables.html

release: release_linux release_windows