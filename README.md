# golang-todo
The simplest to-do app that i use for work \ life notes.

## Usage
If you want to use this app, you shoud change db_path for your SQLite3 DB file. 
```go
db, err = sql.Open("sqlite3", "/home/aaddssww/tasks.db")
```
then you can run it as a regular go programm
```go
go run main.go
```
or build it
```go
go build -o <name_for_your_app> .
```
