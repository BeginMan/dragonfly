# Binary name
BINARY=dragonfly
VERSION="0.1"
# Builds the project
build:
		go build -o ${BINARY} -ldflags "-X main.version=${VERSION}"
		go test -v
# Installs our project: copies binaries
install:
		go install
release:
		# Clean
		go clean
		rm -rf *.gz
		# Build for mac
		go build -ldflags "-X main.version=${VERSION}"
		tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
		go clean

		## Build for arm
		go clean
		CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}

		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=7 CC=arm-linux-gnueabihf-gcc-5 CGO_ENABLED=1 go build -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}

		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}"
		tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
		go clean
# Cleans our projects: deletes binaries
clean:
		go clean
		rm -rf *.gz

.PHONY:  clean build
