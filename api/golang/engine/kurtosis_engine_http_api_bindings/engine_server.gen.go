// Package kurtosis_engine_http_api_bindings provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package kurtosis_engine_http_api_bindings

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Delete Enclaves
	// (DELETE /enclaves)
	DeleteEnclaves(ctx echo.Context, params DeleteEnclavesParams) error
	// Get Enclaves
	// (GET /enclaves)
	GetEnclaves(ctx echo.Context) error
	// Create Enclave
	// (POST /enclaves)
	PostEnclaves(ctx echo.Context) error
	// Get Historical Enclaves
	// (GET /enclaves/historical)
	GetEnclavesHistorical(ctx echo.Context) error
	// Destroy Enclave
	// (DELETE /enclaves/{enclave_identifier})
	DeleteEnclavesEnclaveIdentifier(ctx echo.Context, enclaveIdentifier string) error
	// Get Enclave Info
	// (GET /enclaves/{enclave_identifier})
	GetEnclavesEnclaveIdentifier(ctx echo.Context, enclaveIdentifier string) error
	// Get Service Logs
	// (POST /enclaves/{enclave_identifier}/logs)
	PostEnclavesEnclaveIdentifierLogs(ctx echo.Context, enclaveIdentifier string) error
	// Stop Enclave
	// (POST /enclaves/{enclave_identifier}/stop)
	PostEnclavesEnclaveIdentifierStop(ctx echo.Context, enclaveIdentifier string) error
	// Get Engine Info
	// (GET /engine/info)
	GetEngineInfo(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// DeleteEnclaves converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteEnclaves(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteEnclavesParams
	// ------------- Optional query parameter "remove_all" -------------

	err = runtime.BindQueryParameter("form", true, false, "remove_all", ctx.QueryParams(), &params.RemoveAll)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter remove_all: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteEnclaves(ctx, params)
	return err
}

// GetEnclaves converts echo context to params.
func (w *ServerInterfaceWrapper) GetEnclaves(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEnclaves(ctx)
	return err
}

// PostEnclaves converts echo context to params.
func (w *ServerInterfaceWrapper) PostEnclaves(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostEnclaves(ctx)
	return err
}

// GetEnclavesHistorical converts echo context to params.
func (w *ServerInterfaceWrapper) GetEnclavesHistorical(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEnclavesHistorical(ctx)
	return err
}

// DeleteEnclavesEnclaveIdentifier converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteEnclavesEnclaveIdentifier(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "enclave_identifier" -------------
	var enclaveIdentifier string

	err = runtime.BindStyledParameterWithLocation("simple", false, "enclave_identifier", runtime.ParamLocationPath, ctx.Param("enclave_identifier"), &enclaveIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter enclave_identifier: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteEnclavesEnclaveIdentifier(ctx, enclaveIdentifier)
	return err
}

// GetEnclavesEnclaveIdentifier converts echo context to params.
func (w *ServerInterfaceWrapper) GetEnclavesEnclaveIdentifier(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "enclave_identifier" -------------
	var enclaveIdentifier string

	err = runtime.BindStyledParameterWithLocation("simple", false, "enclave_identifier", runtime.ParamLocationPath, ctx.Param("enclave_identifier"), &enclaveIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter enclave_identifier: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEnclavesEnclaveIdentifier(ctx, enclaveIdentifier)
	return err
}

// PostEnclavesEnclaveIdentifierLogs converts echo context to params.
func (w *ServerInterfaceWrapper) PostEnclavesEnclaveIdentifierLogs(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "enclave_identifier" -------------
	var enclaveIdentifier string

	err = runtime.BindStyledParameterWithLocation("simple", false, "enclave_identifier", runtime.ParamLocationPath, ctx.Param("enclave_identifier"), &enclaveIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter enclave_identifier: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostEnclavesEnclaveIdentifierLogs(ctx, enclaveIdentifier)
	return err
}

// PostEnclavesEnclaveIdentifierStop converts echo context to params.
func (w *ServerInterfaceWrapper) PostEnclavesEnclaveIdentifierStop(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "enclave_identifier" -------------
	var enclaveIdentifier string

	err = runtime.BindStyledParameterWithLocation("simple", false, "enclave_identifier", runtime.ParamLocationPath, ctx.Param("enclave_identifier"), &enclaveIdentifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter enclave_identifier: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostEnclavesEnclaveIdentifierStop(ctx, enclaveIdentifier)
	return err
}

// GetEngineInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetEngineInfo(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEngineInfo(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE(baseURL+"/enclaves", wrapper.DeleteEnclaves)
	router.GET(baseURL+"/enclaves", wrapper.GetEnclaves)
	router.POST(baseURL+"/enclaves", wrapper.PostEnclaves)
	router.GET(baseURL+"/enclaves/historical", wrapper.GetEnclavesHistorical)
	router.DELETE(baseURL+"/enclaves/:enclave_identifier", wrapper.DeleteEnclavesEnclaveIdentifier)
	router.GET(baseURL+"/enclaves/:enclave_identifier", wrapper.GetEnclavesEnclaveIdentifier)
	router.POST(baseURL+"/enclaves/:enclave_identifier/logs", wrapper.PostEnclavesEnclaveIdentifierLogs)
	router.POST(baseURL+"/enclaves/:enclave_identifier/stop", wrapper.PostEnclavesEnclaveIdentifierStop)
	router.GET(baseURL+"/engine/info", wrapper.GetEngineInfo)

}

type DeleteEnclavesRequestObject struct {
	Params DeleteEnclavesParams
}

type DeleteEnclavesResponseObject interface {
	VisitDeleteEnclavesResponse(w http.ResponseWriter) error
}

type DeleteEnclaves200JSONResponse DeleteResponse

func (response DeleteEnclaves200JSONResponse) VisitDeleteEnclavesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetEnclavesRequestObject struct {
}

type GetEnclavesResponseObject interface {
	VisitGetEnclavesResponse(w http.ResponseWriter) error
}

type GetEnclaves200JSONResponse GetEnclavesResponse

func (response GetEnclaves200JSONResponse) VisitGetEnclavesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostEnclavesRequestObject struct {
	Body *PostEnclavesJSONRequestBody
}

type PostEnclavesResponseObject interface {
	VisitPostEnclavesResponse(w http.ResponseWriter) error
}

type PostEnclaves200JSONResponse CreateEnclaveResponse

func (response PostEnclaves200JSONResponse) VisitPostEnclavesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetEnclavesHistoricalRequestObject struct {
}

type GetEnclavesHistoricalResponseObject interface {
	VisitGetEnclavesHistoricalResponse(w http.ResponseWriter) error
}

type GetEnclavesHistorical200JSONResponse GetExistingAndHistoricalEnclaveIdentifiersResponse

func (response GetEnclavesHistorical200JSONResponse) VisitGetEnclavesHistoricalResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type DeleteEnclavesEnclaveIdentifierRequestObject struct {
	EnclaveIdentifier string `json:"enclave_identifier"`
}

type DeleteEnclavesEnclaveIdentifierResponseObject interface {
	VisitDeleteEnclavesEnclaveIdentifierResponse(w http.ResponseWriter) error
}

type DeleteEnclavesEnclaveIdentifier200JSONResponse map[string]interface{}

func (response DeleteEnclavesEnclaveIdentifier200JSONResponse) VisitDeleteEnclavesEnclaveIdentifierResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetEnclavesEnclaveIdentifierRequestObject struct {
	EnclaveIdentifier string `json:"enclave_identifier"`
}

type GetEnclavesEnclaveIdentifierResponseObject interface {
	VisitGetEnclavesEnclaveIdentifierResponse(w http.ResponseWriter) error
}

type GetEnclavesEnclaveIdentifier200JSONResponse EnclaveInfo

func (response GetEnclavesEnclaveIdentifier200JSONResponse) VisitGetEnclavesEnclaveIdentifierResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostEnclavesEnclaveIdentifierLogsRequestObject struct {
	EnclaveIdentifier string `json:"enclave_identifier"`
	Body              *PostEnclavesEnclaveIdentifierLogsJSONRequestBody
}

type PostEnclavesEnclaveIdentifierLogsResponseObject interface {
	VisitPostEnclavesEnclaveIdentifierLogsResponse(w http.ResponseWriter) error
}

type PostEnclavesEnclaveIdentifierLogs200JSONResponse GetServiceLogsResponse

func (response PostEnclavesEnclaveIdentifierLogs200JSONResponse) VisitPostEnclavesEnclaveIdentifierLogsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostEnclavesEnclaveIdentifierStopRequestObject struct {
	EnclaveIdentifier string `json:"enclave_identifier"`
}

type PostEnclavesEnclaveIdentifierStopResponseObject interface {
	VisitPostEnclavesEnclaveIdentifierStopResponse(w http.ResponseWriter) error
}

type PostEnclavesEnclaveIdentifierStop200JSONResponse map[string]interface{}

func (response PostEnclavesEnclaveIdentifierStop200JSONResponse) VisitPostEnclavesEnclaveIdentifierStopResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetEngineInfoRequestObject struct {
}

type GetEngineInfoResponseObject interface {
	VisitGetEngineInfoResponse(w http.ResponseWriter) error
}

type GetEngineInfo200JSONResponse GetEngineInfoResponse

func (response GetEngineInfo200JSONResponse) VisitGetEngineInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Delete Enclaves
	// (DELETE /enclaves)
	DeleteEnclaves(ctx context.Context, request DeleteEnclavesRequestObject) (DeleteEnclavesResponseObject, error)
	// Get Enclaves
	// (GET /enclaves)
	GetEnclaves(ctx context.Context, request GetEnclavesRequestObject) (GetEnclavesResponseObject, error)
	// Create Enclave
	// (POST /enclaves)
	PostEnclaves(ctx context.Context, request PostEnclavesRequestObject) (PostEnclavesResponseObject, error)
	// Get Historical Enclaves
	// (GET /enclaves/historical)
	GetEnclavesHistorical(ctx context.Context, request GetEnclavesHistoricalRequestObject) (GetEnclavesHistoricalResponseObject, error)
	// Destroy Enclave
	// (DELETE /enclaves/{enclave_identifier})
	DeleteEnclavesEnclaveIdentifier(ctx context.Context, request DeleteEnclavesEnclaveIdentifierRequestObject) (DeleteEnclavesEnclaveIdentifierResponseObject, error)
	// Get Enclave Info
	// (GET /enclaves/{enclave_identifier})
	GetEnclavesEnclaveIdentifier(ctx context.Context, request GetEnclavesEnclaveIdentifierRequestObject) (GetEnclavesEnclaveIdentifierResponseObject, error)
	// Get Service Logs
	// (POST /enclaves/{enclave_identifier}/logs)
	PostEnclavesEnclaveIdentifierLogs(ctx context.Context, request PostEnclavesEnclaveIdentifierLogsRequestObject) (PostEnclavesEnclaveIdentifierLogsResponseObject, error)
	// Stop Enclave
	// (POST /enclaves/{enclave_identifier}/stop)
	PostEnclavesEnclaveIdentifierStop(ctx context.Context, request PostEnclavesEnclaveIdentifierStopRequestObject) (PostEnclavesEnclaveIdentifierStopResponseObject, error)
	// Get Engine Info
	// (GET /engine/info)
	GetEngineInfo(ctx context.Context, request GetEngineInfoRequestObject) (GetEngineInfoResponseObject, error)
}

type StrictHandlerFunc func(ctx echo.Context, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// DeleteEnclaves operation middleware
func (sh *strictHandler) DeleteEnclaves(ctx echo.Context, params DeleteEnclavesParams) error {
	var request DeleteEnclavesRequestObject

	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteEnclaves(ctx.Request().Context(), request.(DeleteEnclavesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteEnclaves")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(DeleteEnclavesResponseObject); ok {
		return validResponse.VisitDeleteEnclavesResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetEnclaves operation middleware
func (sh *strictHandler) GetEnclaves(ctx echo.Context) error {
	var request GetEnclavesRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetEnclaves(ctx.Request().Context(), request.(GetEnclavesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetEnclaves")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetEnclavesResponseObject); ok {
		return validResponse.VisitGetEnclavesResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostEnclaves operation middleware
func (sh *strictHandler) PostEnclaves(ctx echo.Context) error {
	var request PostEnclavesRequestObject

	var body PostEnclavesJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostEnclaves(ctx.Request().Context(), request.(PostEnclavesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostEnclaves")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostEnclavesResponseObject); ok {
		return validResponse.VisitPostEnclavesResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetEnclavesHistorical operation middleware
func (sh *strictHandler) GetEnclavesHistorical(ctx echo.Context) error {
	var request GetEnclavesHistoricalRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetEnclavesHistorical(ctx.Request().Context(), request.(GetEnclavesHistoricalRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetEnclavesHistorical")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetEnclavesHistoricalResponseObject); ok {
		return validResponse.VisitGetEnclavesHistoricalResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// DeleteEnclavesEnclaveIdentifier operation middleware
func (sh *strictHandler) DeleteEnclavesEnclaveIdentifier(ctx echo.Context, enclaveIdentifier string) error {
	var request DeleteEnclavesEnclaveIdentifierRequestObject

	request.EnclaveIdentifier = enclaveIdentifier

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteEnclavesEnclaveIdentifier(ctx.Request().Context(), request.(DeleteEnclavesEnclaveIdentifierRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteEnclavesEnclaveIdentifier")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(DeleteEnclavesEnclaveIdentifierResponseObject); ok {
		return validResponse.VisitDeleteEnclavesEnclaveIdentifierResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetEnclavesEnclaveIdentifier operation middleware
func (sh *strictHandler) GetEnclavesEnclaveIdentifier(ctx echo.Context, enclaveIdentifier string) error {
	var request GetEnclavesEnclaveIdentifierRequestObject

	request.EnclaveIdentifier = enclaveIdentifier

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetEnclavesEnclaveIdentifier(ctx.Request().Context(), request.(GetEnclavesEnclaveIdentifierRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetEnclavesEnclaveIdentifier")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetEnclavesEnclaveIdentifierResponseObject); ok {
		return validResponse.VisitGetEnclavesEnclaveIdentifierResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostEnclavesEnclaveIdentifierLogs operation middleware
func (sh *strictHandler) PostEnclavesEnclaveIdentifierLogs(ctx echo.Context, enclaveIdentifier string) error {
	var request PostEnclavesEnclaveIdentifierLogsRequestObject

	request.EnclaveIdentifier = enclaveIdentifier

	var body PostEnclavesEnclaveIdentifierLogsJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostEnclavesEnclaveIdentifierLogs(ctx.Request().Context(), request.(PostEnclavesEnclaveIdentifierLogsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostEnclavesEnclaveIdentifierLogs")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostEnclavesEnclaveIdentifierLogsResponseObject); ok {
		return validResponse.VisitPostEnclavesEnclaveIdentifierLogsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// PostEnclavesEnclaveIdentifierStop operation middleware
func (sh *strictHandler) PostEnclavesEnclaveIdentifierStop(ctx echo.Context, enclaveIdentifier string) error {
	var request PostEnclavesEnclaveIdentifierStopRequestObject

	request.EnclaveIdentifier = enclaveIdentifier

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostEnclavesEnclaveIdentifierStop(ctx.Request().Context(), request.(PostEnclavesEnclaveIdentifierStopRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostEnclavesEnclaveIdentifierStop")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostEnclavesEnclaveIdentifierStopResponseObject); ok {
		return validResponse.VisitPostEnclavesEnclaveIdentifierStopResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetEngineInfo operation middleware
func (sh *strictHandler) GetEngineInfo(ctx echo.Context) error {
	var request GetEngineInfoRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetEngineInfo(ctx.Request().Context(), request.(GetEngineInfoRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetEngineInfo")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetEngineInfoResponseObject); ok {
		return validResponse.VisitGetEngineInfoResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xY34/ithP/VyJ/v4/psW3feKML3UO6BbRkpatOJ8ubTMCnxM7ZE3roxP9e2fkJcSCw",
	"t5Wq9g2S8cx85sdnxvlOQplmUoBATcbfiQ63kDL7c5LxeymQcQFqjQxz+xREnpLxJ/L0vFjMFw/EJ+tg",
	"uVrNpsQni+WCzj7O18FsEZDPPsF9BmRMNCouNuTgk3sFDGEmwoTtYKI2VmOmZAYKOdh/LOM0rMzSRG5o",
	"AjtIzKuOumPZHSjNpaDINk5pKMxSwVJwCqQysi/+ryAmY/K/UROaURmXUen7oxE9HGqI8uULhNiB+AQ6",
	"k0JDF2blDBexHGhzbkSdNqeQAJ4xpiCVO4hoOwKUiYjmOY+sBEdI9UA/FiyFiYiecx6RxhumFNu73avy",
	"vZrX9fReanxk4ZaLAlbH5Y3KQppJhVQKupUaaVqItzLHBcIGlDHBszNyVYYH+uZ26EXxaAOUZ5RFkQKt",
	"nSXUVCOPnAINLi40j6DKSS+sXrEBqGpIelj/zh5XwR/Oxq1KMAKBPOagdH9Jm6JyYu9tPL2VCkFA1Hf2",
	"DEZ3to6ZoV0X13TcuYrt0M+Net3KdJ2vc+ocDN2uQT1QTV+1GF2GzSyp8vQiNwY8BY0szdps21sNV7Pt",
	"GxXQY+lH1RjBbB0Qn6yeltPn+2C+XJxriDYXdoqw190rnHwALE3p4cOERRE3WWPJ6khu8JDpONLr2qbs",
	"h3POGZlqOl8B+xvXyMVmIqL3XKNUPGRJl4f6DbMkOeGra6Zc++igKfcAuAa14yF8kBvtXm5CKb7kIkS+",
	"AxrzBK/x64PcfOACfrfHui75JJZJIv80K1N7Nr1ImQATtn3ytFiouADtHjgKMFeCsiQ5o0cXMG2vUQ14",
	"BKFT7NdGrj+hQiKNZS6M0ds9aPw3COnL/kjbre1TZmdg61TSHYhJ+XQ4GKxJdzg7n3GpLK+OY+Y3Q6kG",
	"xmFZiRtT8A1pxhBBDW3/Uy0tep4uZ2t6v1wEk/mCBrOPhqvts8UycD6vnj1Ogvv39Gn2MPvoOtJ+7SL8",
	"oB3nWKqUIRmTiCH8ZGdj54gBVhEyckzMu4IvvclqTnxSMyL5+d3duztjRGYgWMbJmPxqH/kkY7i1CRiV",
	"JG//RHbhb/LCpZhHZFxeBKqBYY8rlkLBM5/MOR0qnmFhdh57qHLwveJ+4LEk8Sor77wpxCxP0OPai1mi",
	"DUJuTn3NQe1JNYzLu4VhDOKXt0cXaxw+G3IpOttC+OXurmREBIHFzpYlPLRgRl90MSoaheeK7uT+Y0N/",
	"DHWdhyFoHeeJpxoxn+g8TZna16Hz6tiZXb0gluMQtwYyeUNIrrl/E64HwCNQmdQOVCupj2F9zUHjbzLa",
	"/zBE3bu/xWNMcQURGZtaPLxhSN0385uCWqiq4mpf1v052tbbir3IXiiiZrd563K6cp+6udoa7a3CO4rQ",
	"93phre0ehtNax+dLPPf8PJ/6Xn1H8Ir/UnmGwjwZe7iFivgqljO825Bc119yWrkO8qsnwWu573SjuIne",
	"NCq5b0r2Mr39++I8/CPc64jYqz/lne+JUbWDXybtTrLMKv3PSNiPHzWOq9jfPGt6rjQ3F06pzLNZHVA4",
	"GmV2Y+GszdH/GPViWkycTjcAs92PqqW/n1+rjyZvvkB2vs68grnszaWiv8NfAQAA//+kzK+LOBoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}