package integration_test

import (
	. "github.com/Eun/go-hit"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	attempts   = 20
	host       = "app:4000"
	healthPath = "http://" + host + "/healthcheck"

	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

func TestHTTPCreateDevice(t *testing.T) {
	body := `{
		"name": "foo",
		"value": 2221111,
		"description": "bar"
	}`
	Test(t,
		Description("CreateDevice Success"),
		Post(basePath+"/devices"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),

		Expect().Status().Equal(http.StatusCreated),
		Expect().Headers("Location").Len().GreaterThan(0),
		Expect().Body().JSON().JQ(".device.name").Equal("foo"),
		Expect().Body().JSON().JQ(".device.value").Equal(2221111),
		Expect().Body().JSON().JQ(".device.description").Equal("bar"),
	)

	body = `{
		"name": "",
		"value": 2221111,
		"description": "bar"
	}`

	Test(t,
		Description("CreateDevice Fail"),
		Post(basePath+"/devices"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),

		Expect().Status().Equal(http.StatusUnprocessableEntity),
		Expect().Body().JSON().JQ(".error.name").Equal("must be provided"),
	)
}

func TestHTTPShowDevice(t *testing.T) {
	Test(t,
		Description("ShowDevice Success"),
		Get(basePath+"/devices/1"),
		Send().Headers("Content-Type").Add("application/json"),

		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".device.id").Equal(1),
		Expect().Body().JSON().JQ(".device.name").Equal("foo"),
		Expect().Body().JSON().JQ(".device.value").Equal(2221111),
		Expect().Body().JSON().JQ(".device.description").Equal("bar"),
	)

	Test(t,
		Description("ShowDevice Fail"),
		Get(basePath+"/devices/n"),

		Expect().Status().Equal(http.StatusNotFound),
		Expect().Body().String().Contains("the requested resource could not be found"),
	)
}
