package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/TylerBrock/colorjson"
)

func convertHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, value := range headers {
		result[key] = value[0]
	}
	return result
}

func convertBody(body io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func convertParams(params map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, value := range params {
		result[key] = value[0]
	}
	return result
}

func checkfIfValidIp(ip string) bool {
	return !(net.ParseIP(ip) == nil)
}

func checkIfValidPort(bind string) bool {
	port, err := strconv.Atoi(bind)
	if err != nil {
		return false
	}
	return port > 0 && port < 65535
}

func main() {
	bind := ""
	flag.StringVar(&bind, "bind", bind, "server bind address (with port)")
	flag.StringVar(&bind, "b", bind, "server bind address (with port)")
	flag.Parse()

	splitted := strings.Split(bind, ":")
	if len(splitted) != 2 {
		if checkfIfValidIp(bind) {
			bind = bind + ":8080"
		} else if checkIfValidPort(bind) {
			bind = ":" + bind
		} else {
			log.Println("Invalid bind address value:", bind)
			os.Exit(1)
		}
	}

	e := echo.New()
	e.Any("*", func(c echo.Context) error {
		formParams, _ := c.FormParams()

		sendData := map[string]interface{}{
			"method":      c.Request().Method,
			"uri":         c.Request().RequestURI,
			"headers":     convertHeaders(c.Request().Header),
			"body":        convertBody(c.Request().Body),
			"queryParams": convertParams(c.QueryParams()),
			"formParams":  convertParams(formParams),
			"remoteAddr":  c.Request().RemoteAddr,
			"host":        c.Request().Host,
			"protocol":    c.Request().Proto,
			"referer":     c.Request().Referer(),
			"reguestTime": time.Now().Format(time.RFC3339),
		}

		f := colorjson.NewFormatter()
		f.Indent = 4

		s, _ := f.Marshal(sendData)
		fmt.Println(string(s))

		return c.JSON(http.StatusOK, sendData)
	})
	e.Logger.Fatal(e.Start(bind))
}
