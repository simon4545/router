GOOS=linux go build -ldflags '-w -s' ./main.go
upx -9 main