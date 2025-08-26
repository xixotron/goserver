# Simple http Go server

This is a simple HttpServer in Go

This server is being built as I progress through [boot.dev](https://boot.dev)'s
course [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang).

## Functionality

So far we start a server on port _8080_
now we serve the files in the folder we run the server from /app path

We now have an API /api with the method:

/api/healthz  allways returns 200/OK when the server is runing

Now we have metrics under /admin directory:

- /admin/metrics  tells us how many times the /app has ben called (Hits)
- /admin/reset    resets the number of Hits on the /app counter
