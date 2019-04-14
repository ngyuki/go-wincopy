
all: wincopy.exe

wincopy.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o wincopy.exe .
