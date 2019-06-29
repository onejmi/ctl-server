# CTL Server
Backend for Claim The Light project.

Written in Go.

### TODO
* ~~Add profile route (GET) which will serve user objects to authenticated clients~~
* Limit the amount of requests a certain session can send to the server in a certain period of time (rate limiting)
* Track user IPs
* ~~Refactor error "code" to "status"~~

### How to run
Use `go run claim_the_light.go` assuming Go is apart of the `$PATH` environment variable
and the command is being executed in the **src/** directory.