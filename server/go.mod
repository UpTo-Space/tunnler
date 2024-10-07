module github.com/UpTo-Space/tunnler/server

go 1.23.2

require (
	github.com/UpTo-Space/tunnler/common v0.0.0
	github.com/coder/websocket v1.8.12
	github.com/golang-migrate/migrate/v4 v4.18.1
	github.com/lib/pq v1.10.9
	golang.org/x/time v0.5.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)

replace github.com/UpTo-Space/tunnler/common v0.0.0 => ../common
