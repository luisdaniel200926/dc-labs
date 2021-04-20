User Guide
==========

To setup the programm you will need to instal Golang, Gin Package, jwt-go 

https://golang.org/doc/install

https://github.com/gin-gonic/gin#installation

https://github.com/dgrijalva/jwt-go


Install and run the server using this command:	go run server.go

The input in console should be as it follows:

To login:
	curl -u username:password http://localhost:8080/login

To logout:
	curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout

To make an upload:
	curl -F 'data=@/path/to/local/image.png' -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/upload

To get status:
	curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status
