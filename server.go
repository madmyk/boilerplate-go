// Package main is a boilerplate for JSON based microservices.
// This documentation should reflect the scope of the service.
//
// Microservices are a collection of loosely coupled services,
// all functionality within this single source file should
// reflect that.
//
// Encapsulate complex business logic into shareable libraries
// abstracted away from the service implementation expressed
// in the source below.
//
// Simple business logic may be implemented here if it adds
// little or no significant complexity or limits readability
// or the surrounding source
//
// Compile with -tags=jsoniter for faster json performance.
// https://github.com/gin-gonic/gin#build-with-jsoniter
//
// txn2.com
package main

import (
	"os"
	"time"

	"io/ioutil"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	// Default and consistent environment variables
	// help standardize k8s configs and documentation
	//
	port := getEnv("PORT", "8080")
	debug := getEnv("DEBUG", "false")
	cfgFile := getEnv("CONFIG", "")
	basePath := getEnv("BASE_PATH", "")

	// load a configuration yml is one is specified
	//
	cfg := make(map[interface{}]interface{})
	if cfgFile != "" {
		ymlData, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal([]byte(ymlData), &cfg)
		if err != nil {
			panic(err)
		}
	}

	gin.SetMode(gin.ReleaseMode)

	if debug == "true" {
		gin.SetMode(gin.DebugMode)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err.Error())
	}

	if debug == "true" {
		logger, _ = zap.NewDevelopment()
	}

	// router
	r := gin.New()
	rg := r.Group(basePath)

	// middleware
	//
	rg.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// routes
	//
	rg.GET("/",
		func(c *gin.Context) {

			// call external libs for business logic here

			// return
			c.JSON(200, gin.H{"status": "success"})
			return
		},
	)

	rg.POST("/",
		func(c *gin.Context) {
			rs, err := c.GetRawData()
			if err != nil {
				c.JSON(500, gin.H{"status": "fail", "error": err.Error()})
				return
			}
			// parse json validation etc..
			// call external libs for business logic here

			c.JSON(200, gin.H{"status": "success", "body": string(rs)})
			return
		},
	)

	// for external status check
	rg.GET("/status",
		func(c *gin.Context) {
			// return
			c.JSON(200, gin.H{
				"status": "success",
				// add additional metrics here
			})
			return
		},
	)

	// default no route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": "fail", "code": 404, "message": "not found"})
	})

	r.Run(":" + port)
}

// getEnv gets an environment variable or sets a default if
// one does not exist.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
