package main

import (
	"fmt"
	"os"
	"io"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string
	Password string
	Token string
}

var loggedUsers []*User

/*
A sample user
*/

var defaultUser = User {
	Username: "username",
	Password: "password",
}

// new user recives the username and the password and returns a user struct
func newUser(username string, password string) *User{
	token, _ := createToken(username)				// handle error

	u := User{
		Username: username,
		Password: password,
		Token: token,
	}
	return &u
}

// recives the username and returns a token (string) and an error if exists, if not nil
func createToken(username string) (string, error) {
	var err error

	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")			// ad to env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

// login API endpoint
func login(c *gin.Context) {
	username, password, _ := c.Request.BasicAuth()
	if username != defaultUser.Username || password != defaultUser.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid credentials")
		return
	}

	user := newUser(username, password)
	loggedUsers = append(loggedUsers, user)

	message := fmt.Sprintf("Hi %s welcome to the DPIP system", user.Username)

	c.JSON(200, gin.H {
		"message": message,
		"token": user.Token,
	})
}

// struct tu recieve the token from the header
type tokenHeader struct {
	Token string `header:"Authorization"`
}

// logout API endpoint
func logout(c *gin.Context) {
	h := tokenHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	for i, u := range loggedUsers {
		if h.Token == u.Token {
			message := fmt.Sprintf("Bye %s, your token has been revoked", u.Username)

			loggedUsers = append(loggedUsers[:i], loggedUsers[i + 1:]...)

			c.JSON(200, gin.H {
				"message": message,
			})
			return
		}

		c.JSON(500, gin.H {					// err 500 -> bad token
			"message": "Error, not logged in",
		})
	}
}

// struct to recieve the data (file to upload path) from the body
type uploadBody struct {
	Body string `content-type: "JSON" form:"data"`
}

// upload API endpoint
func upload(c *gin.Context) {
	h := tokenHeader{}
	//data := uploadBody{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
		return
	}

	for _, u := range loggedUsers {
		if h.Token == u.Token {
			multipartFile, err := c.FormFile("data")

			if err != nil {
				c.JSON(100, gin.H {				// 100 -> file upload error
					"message": err,
				})
				return
			}

			f, err := multipartFile.Open()

			if err != nil {
				c.JSON(104, gin.H {
					"message": err,				// 104 -> file open error
				})
				return
			}

			newFile, err := os.Create("./imgs/" + multipartFile.Filename)

			if err != nil {
				c.JSON(101, gin.H {				// 101 -> file create error
					"message": err,
				})
				return
			}

			if _, err := io.Copy(newFile, f); err != nil {
				c.JSON(102, gin.H {				// 102 -> file copy error
					"message": err,
				})
				return

			}

			fi, err := newFile.Stat()
			if err != nil {
				c.JSON(103, gin.H {				// 103 -> file size error
					"message": err,
				})
				return
			}

			size := fmt.Sprintf("%d bytes", fi.Size())

			c.JSON(200, gin.H {
				"message": "An image has been successfully uploaded",
				"filename": multipartFile.Filename,
				"size": size,
			})

			f.Close()
			newFile.Close()
			return
		}
	}

	c.JSON(500, gin.H {					// err 500 -> bad token
		"message": "Error, not logged in",
	})
}

// status API enpoint
func status(c *gin.Context) {
	h := tokenHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	for _, u := range loggedUsers {
		if h.Token == u.Token {
			message := fmt.Sprintf("Hi %s, the DPIP System is Up and Running", u.Username)

			c.JSON(200, gin.H {
				"message": message,
				"time": time.Now().Format("2006-Jan-02 15:04:05"),//("2015-03-07 11:06:39"),
			})
			return
		}
	}

	c.JSON(500, gin.H {					// err 500 -> bad token
		"message": "Error, not logged in",
	})
}

// main function
func main() {
	fmt.Printf("Starting server at port 8080\n")

	server := gin.Default()

	server.GET("/login", login)
	server.GET("/logout", logout)
	server.GET("/status", status)

	server.MaxMultipartMemory = 8 << 20			// 8 MiB
	server.POST("/upload", upload)

	server.Run()
}

/*
errors

err 100 -> file upload error
err 101 -> file create error
err 102 -> file copy error
err 103 -> file size error
err 104 -> file open error
err 500 -> bad token
err 700 -> header error


*/
