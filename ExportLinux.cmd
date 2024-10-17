if exist TitanAttendance del /Q TitanAttendance
go get -u
go mod tidy
set GOOS=linux
set GOARCH=amd64
go build -o TitanAttendance -trimpath -ldflags="-s -w"