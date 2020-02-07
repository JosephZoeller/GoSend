#!/bin/bash

cd ..
docker-compose up --no-start

xterm -xrm 'XTerm.vt100.allowTitleOps: false' -T "<REVERSE PROXY>" -e docker start -a rvprox &
xterm -xrm 'XTerm.vt100.allowTitleOps: false' -T "<LOG MANAGER>" -e docker start -a logmgr &
xterm -xrm 'XTerm.vt100.allowTitleOps: false' -T "<SERVER 1>" -e docker start -a srvr_1 &
xterm -xrm 'XTerm.vt100.allowTitleOps: false' -T "<SERVER 2>" -e docker start -a srvr_2 &

cd test
x-terminal-emulator
# ./client -out=:8080 -files=test.jpeg