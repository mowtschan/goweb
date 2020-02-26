package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	envPort       = "PORT"
	envOutputFile = "OUTPUT_FILE"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetOutput(os.Stderr)

	e.POST("/", handler)
	port := os.Getenv(envPort)
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func handler(c echo.Context) error {
	data := new(map[string]interface{})
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	outputFile := os.Getenv(os.Getenv(envOutputFile))
	if outputFile == "" {
		outputFile = "/tmp/payload.json"
	}
	content := new(bytes.Buffer)
	if err := json.NewEncoder(content).Encode(data); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if err := ioutil.WriteFile(outputFile, content.Bytes(), 0777); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}
