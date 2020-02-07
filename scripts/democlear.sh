#!/bin/bash
cd ..
docker-compose stop
docker rm rvprox logmgr srvr_1 srvr_2
docker rmi logmgr rvprox srvr
docker network rm gmg_default