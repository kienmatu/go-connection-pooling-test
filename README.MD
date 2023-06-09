# Test benchmark connection pooling in Go with Postgres
#### Author: kienmatu

## Introduction
This is a simple project to test benchmark connection pooling in Go with Postgres. 
The purpose of this project is to compare the performance of three different methods of handling database connections: 
- Using connection pooling
- Creating new connections each time
- And using a single connection.

### Installation and Usage:

To run this project, you will need to have Go and docker installed on your system.
Once you have installed these dependencies, you can follow the steps below:

Clone the repository:

```bash
git clone https://github.com/kienmatu/go-connection-pooling-test.git
cd go-connection-pooling-benchmark
```

Install the required dependencies:

```bash
go mod download
```
Start the Postgres server using Docker Compose:

```bash
docker-compose up -d
```
Wait for 1 minute for the data to be seeded.

Run the Sample API:
```bash
go run main.go
```

Install Go-wrk
```bash
go install github.com/tsliwowicz/go-wrk@latest
```
### Test Benchmark
Now you can run the benchmark tests separately using go-wrk:

```bash
go-wrk -c 10 http://localhost:8080/products/pooled
go-wrk -c 10 http://localhost:8080/products/normal
go-wrk -c 10 http://localhost:8080/products/new
```

The results of each test will be displayed in the terminal.

### Test elapsed time
Or you can test it with the `index.html` file:

![test.png](test.png)

### How many connections is the best fit for pool size?
Actually, I set 50 and 90 for the first time bench, after that I realize that too much number of connections is not efficient.

The efficient number of connection pool size should be the same as your server/machine's physical number of cores.

I tested with 4 for both idle connections and max pool size, and the result nearly the same with 50 idle and 90 max conns.


I just also added the code in expressjs (but yeah, of course, the performance can not be strong as Go).