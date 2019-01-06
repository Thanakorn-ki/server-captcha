package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/salapao2136/captcha"
)

func main() {
	uri := fmt.Sprintf("mongodb://root:root@127.0.0.1:27017")

	fmt.Println("Program Start...")
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	collection := session.DB("workshop").C("captcha")

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/captcha", func(c *gin.Context) {
		answer, captcha := generateCaptcha()

		_, err := collection.Upsert(
			bson.M{
				"captcha": captcha,
			},
			bson.M{
				"captcha": captcha,
				"answer":  answer,
			})

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		c.JSON(200, gin.H{
			"message": captcha,
			"answer":  answer,
		})
	})

	r.POST("/captcha", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "jwt",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func generateCaptcha() (int, string) {
	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(8) + 1
	second := rand.Intn(8) + 1
	ops := rand.Intn(2) + 1
	var sum = 0
	if ops == 1 {
		sum = first + second
	} else if ops == 2 {
		sum = first - second
	} else if ops == 3 {
		sum = first * second
	}
	captcha := captcha.Captcha(rand.Intn(1)+1, first, ops, second)
	return sum, captcha
}
