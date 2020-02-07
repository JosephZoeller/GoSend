#!/bin/bash

cd ..
make
docker-compose up --no-start

xterm -T "<REVERSE PROXY>" -e docker start -a rvprox &
xterm -T "<LOG MANAGER>" -e docker start -a logmgr &
xterm -T "<SERVER 1>" -e docker start -a srvr_1 &
xterm -T "<SERVER 2>" -e docker start -a srvr_2 &

cd test
x-terminal-emulator
# ./client -out=:8080 -files=test.jpeg