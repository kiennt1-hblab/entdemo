Framework
---
To install Gin package, you need to install Go and set your Go workspace first.
You first need Go installed (version 1.16+ is required), then you can use the below Go command to install Gin.
```
go get -u github.com/gin-gonic/gin
```
Auto Reload
---

Live reload for Go apps.
Via go install.
With go 1.16 or higher:
```
go install github.com/cosmtrek/air@latest
```
to run in terminal
```
air
```
Database
---

Run the following Go commands to install Ent
```
go get -d entgo.io/ent/cmd/ent
```
Create Your First Schema
```
go run -mod=mod entgo.io/ent/cmd/ent init User
```
Let's run the command for generating assets to interact with the `User` entity:
```
go generate ./ent
```
run every time an entity is edited

Need
---

sql driver: 
```
go get -u github.com/go-sql-driver/mysql
```