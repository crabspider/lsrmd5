all: lsrmd5.exe lsrmd5

lsrmd5.exe: main.go
	GOOS=windows GOARCH=amd64 go build -o lsrmd5.exe main.go

lsrmd5: main.go
	GOOS=darwin GOARCH=arm64 go build -o lsrmd5 main.go

clean:
	rm -f lsrmd5.exe lsrmd5
