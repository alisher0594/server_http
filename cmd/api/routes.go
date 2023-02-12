package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/devices", app.createDeviceHandler)
	router.HandlerFunc(http.MethodGet, "/v1/devices/:id", app.showDeviceHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/devices/:id", app.updateDeviceHandler)
	router.HandlerFunc(http.MethodGet, "/v1/devices", app.listDeviceHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/devices/:id", app.deleteDeviceHandler)

	return app.recoverPanic(app.rateLimit(router))
}
