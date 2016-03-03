#Developing Lark

To develop Lark you will need a POSIX shell environment and the Go compiler.
Additionally, the [glide](https://github.com/Masterminds/glide) tool for vendor
package management is required.

It is easiest to start developing lark if you have already installed a Lark
executable [binary](https://github.com/bmatsuo/lark/releases).  Then the
project lark scripts can be used to install dependencies.

```sh
lark init gen test install
```

If you do not want to install Lark as a binary first you will have to use `go
get` to bootstrap a development environment.

```sh
go get ./cmd/...
./lark init gen test install
```
