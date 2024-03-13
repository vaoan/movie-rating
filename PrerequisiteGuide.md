# Prerequisite Guide

This guide will walk you through how to install and verify the required tools for the interview coding exercise. 
It will also confirm that you can run the exercise applications locally.
Please complete this guide prior to the interview. Total completion time is usually no more than one hour. This can vary based on how many tools are already installed on your machine.

You'll need to be able to run a golang API and a react app on your local machine in order to complete the coding exercise.

Formatted `text like this` represents a command to run in your terminal.

## Install Tools/Languages

### 1) Install your preferred code editor or IDE
The following are list of good ones but feel free to use whatever you would like:

- [Webstorm](https://www.jetbrains.com/webstorm/)
- [GoLand](https://www.jetbrains.com/go/)
- [Visual Studio Code](https://code.visualstudio.com)

### 2) Install NPM
#### Mac
1. Install Homebrew: https://docs.brew.sh/Installation
2. `brew update`
3. `brew doctor`
4. `export PATH="/usr/local/bin:$PATH"` _This adds Homebrew’s location to your $PATH in your .bash_profile or .zshrc file and may not be needed._
5. Install Node: `brew install node`

#### Windows and alternatively Mac
1. https://nodejs.org/en/download/
2. Accept "Automatically install necessary tools"

- Note: Alternatively you can install homebrew for windows and `brew install node` but node installation is easier for windows than brew

### 3) Install Docker
1. Docker installation: https://docs.docker.com/get-docker/

2. **For Windows Machines Only:** To run docker you may need to install the linux kernel update package in step 4: https://docs.microsoft.com/en-us/windows/wsl/install-manual#step-4---download-the-linux-kernel-update-package

### 4) Install go
1. Go installation https://go.dev/doc/install
2. Validate by typing `go version` in your terminal

- Note: Alternatively  you can install with homebrew with `brew install go`


## Start the Go API and React Frontend Locally

### Go API
1. Start docker by opening the docker desktop application.
2. From a terminal opened to this repo `cd api` and run `go mod vendor`
3. `cd ../playground`
4. your path (`pwd`) should be **{system directories}/interview-pre-req-check/playground**
5. run `./build.sh`
  - you may need to either `chmod +x ./build.sh` to make it executable or just run the below commands as an alternative
    ```
      docker-compose down -v --remove-orphans
      docker-compose rm -f -s
      docker-compose up --always-recreate-deps --remove-orphans --renew-anon-volumes --build
      ```
- Postgres db will start
- API will start at localhost:8080
  - note: you may see a "connection refused" until postgres fully stands up 
6. Validate api started correctly by navigating to `http://localhost:8080/api/health` in a browser or run `curl http://localhost:8080/api/health` and confirming response body of **{"health":"OK"}**

Troubleshooting:
- If you encounter any issues building the application before start, try deleting the provided vendor file at /api/vendor and running `go mod tidy` and `go mod vendor`

### React App
1. Open a new terminal in the repo directory and `cd frontend`
- your path (`pwd`) should be **{system directories}/interview-pre-req-check/frontend**
2. `npm install --location=global react-scripts`
3. `npm install`
4. `npm start`
- React app will start at localhost:3000
- Validate react app is working properly by seeing it state in terminal it has started and validate it is talking to your api by navigating to `http://localhost:3000/health` and seeing OK in the browser console logs

Troubleshooting:
- If you encounter any issues with `npm install`, try deleting the node_module directory and package-lock.json file in /frontend and running the commands again.