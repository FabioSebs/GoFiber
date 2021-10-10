# :smiling_imp: :fire: JWT Authentication With GoFiber :fire: :smiling_imp:
>Explained by this video: https://www.youtube.com/watch?v=d4Y2DkKbxM0&t=2684s

Packages | Repository
-------- | -----
GoFiber  | github.com/gofiber/fiber/v2
GORM | gorm.io/gorm
GORM Driver MySQL | gorm.io/driver/mysql
Go Cryptography | golang.org/x/crypto
Jwt-Go | github.com/dgrijalva/jwt-go

## Description: 
This project shows the way GoFiber handles HTTP Requests to an API and authenticates users. Authentication will be done using JWT. Click to learn more about [JWT](https://jwt.io/introduction)

## Folder Structure
* Backend
    * /controllers
        * authController.go
    * /database
        * connection.go
    * /models
        * user.go
    * /routes
        * routes.go
    * main.go

## Disclaimer
> If you wish to follow along this repository I hope that you are familiar with installing packages using go mod init {module}. Additionally are familiar with using Postman or Visual Studio to do HTTP Requests to an endpoint. Also a good amount of knowledge about Go wouldn't hurt either.


## Making a Model for the User

```go
package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-" `
}
```
#### Here we make a simple struct modeling the properties our User. The backtick literals next to the property types are called [Struct Tags](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go). In this case the struct tags are used for encoding how the json keys will be displayed.

## Connecting to a MYSQL Database

```go
package database

import (
	"github.com/FabioSebs/GoFiber/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root@/table_name"), &gorm.Config{})

	if err != nil {
		panic("could not connect to database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
```
#### [GORM]("https://gorm.io/docs/index.html) is an ORM library that lets us connect to a database, do auto migrations, query data within tables, all without any SQL! The Connect() function establishes a connection to the database and simultaneously creates a schema. 

## Creating Routes

```go
package routes

import (
	"github.com/FabioSebs/GoFiber/backend/controllers"
	fiber "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}

```

#### The Setup() function requires a \*fiber.App type to be passed in the arguement which will be passed in *main.go*. Next is specifying the HTTP Requests the server will handle. Lastly within request an endpoint has to be specified, and then a callback function is the next arguement to handle the logic for the request. This will be done in *authControllers.go*

## Handling Requests

```go
package controllers

import (
	"strconv"
	"time"
	"github.com/FabioSebs/GoFiber/backend/database"
	"github.com/FabioSebs/GoFiber/backend/models"
	"github.com/dgrijalva/jwt-go"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	database.DB.Create(&user)

	return c.JSON(user)
}
```

#### Lets begin with the Register() Function. This will accept a \*fiber.Ctx object as an arguement. Next we make a map of string key, value pairs. Next to understand this conditional you'll have to know how pointers work. With Go's feature of pointers , we can check if there is an error parsing the data from the context property of fiber (\*fiber.Ctx) and also fill the map (data) simultaneously. Next up we use the *bcrypt* package to encrypt the password so it's not displayed as a raw string in the database. Finally we make our user model and give it the proper fields as the struct we made earlier follows. The next lines create the user table and returns the user in JSON format as a response.



```go
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour + 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "couldn't login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
```