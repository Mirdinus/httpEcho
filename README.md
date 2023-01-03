# Http echo server

This server returns any request made to it. Always returns status code 200 and JSON object with request's info


### Build
Run `go build . -o httpEcho`

### Run
Build and then run  `./httpEcho`

Possible flags are **b** and **bind**, where **b** is just shorter version of **bind**

Flag **bind** is optional. Values can be for example `8080`, `:8080`, `127.0.0.1:8080`, `127.0.0.1`

*If user specifies only IP address, then default port 8080 will be used*