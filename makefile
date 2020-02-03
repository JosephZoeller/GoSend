all: test_linux

test_linux:
	go build -o ./Test/client ./cmd/client
	go build -o ./Test/rvprox ./cmd/rvprox
	go build -o ./Test/server1/server ./cmd/server
	go build -o ./Test/server2/server ./cmd/server
	cp ./scripts/clienttest.sh ./Test/clienttest.sh
	cp ./scripts/networkinittest.sh ./Test/networkinittest.sh

release_linux:
	go build -o ./Release_Build/linux-amd64/client ./cmd/client
	go build -o ./Release_Build/linux-amd64/reverseProxy/GMG-ReverseProxy ./cmd/rvprox
	go build -o ./Release_Build/linux-amd64/server/GMG-Server ./cmd/server
	#mkdir -p ./Release_Build/linux-amd64/server/web
	#cp ./web/tables.html ./Release_Build/linux-amd64/server/web/tables.html

release_windows:
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/client.exe ./cmd/client
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/reverseProxy/GMG-ReverseProxy.exe ./cmd/rvprox
	env GOOS=windows GOARCH=amd64 go build -o ./Release_Build/windows-amd64/server/GMG-Server.exe ./cmd/server
	#mkdir -p ./Release_Build/windows-amd64/server/web/
	#cp ./web/tables.html ./Release_Build/windows-amd64/server/web/tables.html

release: release_linux release_windows