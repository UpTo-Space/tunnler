module github.com/UpTo-Space/tunnler/client

go 1.23.2

require (
	github.com/UpTo-Space/tunnler/common v0.0.0
	github.com/coder/websocket v1.8.12
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/chacha20poly1305 v0.0.0-20170617001512-233f39982aeb // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/o1egl/paseto v1.0.0 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)

replace github.com/UpTo-Space/tunnler/common v0.0.0 => ../common
