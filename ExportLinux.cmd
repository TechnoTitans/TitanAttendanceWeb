if exist TitanAttendance del /Q TitanAttendance
set GOOS=linux
set GOARCH=amd64
go build -o TitanAttendance -trimpath -ldflags="-s -w"