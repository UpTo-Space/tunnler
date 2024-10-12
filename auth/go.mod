module github.com/UpTo-Space/tunnler/auth

go 1.22.1

require github.com/UpTo-Space/tunnler/database v0.0.0

require (
	github.com/golang-migrate/migrate/v4 v4.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)

replace github.com/UpTo-Space/tunnler/database v0.0.0 => ../database
