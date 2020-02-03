#!/bin/bash

x-terminal-emulator -e ./server -in ":8080 :8081"
x-terminal-emulator -e ./server -in ":8082 :8083"

x-terminal-emulator -e ./rvprox -in ":8084 :8085" -out ":8080 :8081 :8082 :8083"

x-terminal-emulator -e ./client -out ":8084" -files "./test.zip"
x-terminal-emulator -e ./client -out ":8085" -files "./test.txt ./test.jpg"