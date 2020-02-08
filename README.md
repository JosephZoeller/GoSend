# GoSend - Client-to-Host TCP File Transfer
This project is a suite of applications for Client-to-Host TCP File Transfer, featuring a reverse proxy with round-robin load balancing and a logging manager. The project includes Docker containerization, with the requisite Dockerfiles and a docker-compose demonstration environment. Miscellaneous testing, deployment and demonstration scripts are available as well.

### Installation
Install this go package with `go get -u github.com/JosephZoeller/gmg`. With docker-compose installed, run the `demo.sh` to launch the default services described in the `docker-compose.yml`. The client application is not containerized in this demo, and should be launched through a separate terminal window.

## Found a bug?

Please submit a bug report to GitHub with as much detail as possible. Please include the log files if applicable.
