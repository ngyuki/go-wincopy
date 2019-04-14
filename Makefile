
all: wincopy.exe

wincopy.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o wincopy.exe .

release: wincopy.exe
	github-release release \
		--user ngyuki \
		--repo go-wincopy \
		--tag $$(git describe --tags) \
		--pre-release
	github-release upload \
		--user ngyuki \
		--repo go-wincopy \
		--tag $$(git describe --tags) \
		--name wincopy.exe \
		--file wincopy.exe \
		--replace
