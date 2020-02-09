#!/bin/bash
cd ..
make
cd test
x-terminal-emulator -e ./logmgr -in ":8086 :8087 :8088" &

x-terminal-emulator -e ./server -in ":8080 :8081" -log ":8086" &
x-terminal-emulator -e ./server -in ":8082 :8083" -log ":8087" &

x-terminal-emulator -e ./rvprox -in ":8084 :8085" -out ":8080 :8081 :8082 :8083" -log ":8088" &

# Assumes test.zip, test.txt and test.jpeg are stored in the /test/ directory
x-terminal-emulator -e ./client -out ":8084" -files ./test.zip &
x-terminal-emulator -e ./client -out ":8085" -files "./test.txt ./test.jpeg"