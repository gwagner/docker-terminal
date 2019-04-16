module Docker-Terminal

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/blinkist/go-dockerpty v0.0.0-20180709141008-29c44b050eff
	github.com/fsouza/go-dockerclient v1.3.6
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/sirupsen/logrus v1.4.1 // indirect
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c // indirect
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
	golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
)

replace github.com/docker/docker v0.0.0-20170601211448-f5ec1e2936dc => github.com/docker/engine v0.0.0-20190408150954-50ebe4562dfc
