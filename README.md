# GMG - Guard My Go
This project is a file transfer application, managed by a Reverse Proxy. Additional features include a round-robin load balancer and a logging manager. The project is containerizable for Docker, containing the requisite Dockerfiles and a docker-compose.yml. Miscellaneous testing, deployment and demonstration scripts are included as well.

### Installation
Install this go package with `go get -u github.com/JosephZoeller/gmg`. With docker-compose installed, run the `demo.sh` to launch the default services described in the `docker-compose.yml`. The client application is not containerized in this demo, and should be launched through a separate terminal window.

## Functions
- [x] Reverse Proxy
- [ ] Firewall
- [ ] Intrusion Detection System
- [x] Logging Manager
- [x] Load Balancer

## Requirements
- [x] Documentation
- [ ] Unit Testing
- [x] Logs & Metrics
- [x] Environment Configuration
- [ ] Security
- [x] Build & Deploy Scripts
- [x] Containerization

## Presentation
- [ ] 10-minute Demonstration
- [ ] Presentation Slides

## Found a bug?

Please submit a bug report to GitHub with as much detail as possible. Please include the log files if applicable.
