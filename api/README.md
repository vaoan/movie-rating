# Movie Rating API

## Getting Started

### Running the Go API Locally
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