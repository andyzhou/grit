## about
this is a simple data storage library, base on level db

## feature
- support simple doc get and save
- support simple and hashed counter get and save

## how to use?
pls see main.go in `example` dir

# Testing
go test -v
go test -bench="Doc" -benchtime=5s .
go test -bench="Count" -benchtime=5s .