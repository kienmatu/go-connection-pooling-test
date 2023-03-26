# go install github.com/tsliwowicz/go-wrk@latest
# run separately
go-wrk -c 10  http://localhost:8080/products/pooled &&
go-wrk -c 10  http://localhost:8080/products/normal

#   Usage: go-wrk <options> <url>
#   Options:
#    -H       Header to add to each request (you can define multiple -H flags) (Default )
#    -M       HTTP method (Default GET)
#    -T       Socket/request timeout in ms (Default 1000)
#    -body    request body string or @filename (Default )
#    -c       Number of goroutines to use (concurrent connections) (Default 10)
#    -ca      CA file to verify peer against (SSL/TLS) (Default )
#    -cert    CA certificate file to verify peer against (SSL/TLS) (Default )
#    -d       Duration of test in seconds (Default 10)
#    -f       Playback file name (Default <empty>)
#    -help    Print help (Default false)
#    -host    Host Header (Default )
#    -http    Use HTTP/2 (Default true)
#    -key     Private key file name (SSL/TLS (Default )
#    -no-c    Disable Compression - Prevents sending the "Accept-Encoding: gzip" header (Default false)
#    -no-ka   Disable KeepAlive - prevents re-use of TCP connections between different HTTP requests (Default false)
#    -no-vr   Skip verifying SSL certificate of the server (Default false)
#    -redir   Allow Redirects (Default false)
#    -v       Print version details (Default false)