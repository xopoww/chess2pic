// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/xopoww/chess2pic/internal/chess2pic"
	"github.com/xopoww/chess2pic/models"
	"github.com/xopoww/chess2pic/pkg/chess"
	"github.com/xopoww/chess2pic/pkg/pic"
	"github.com/xopoww/chess2pic/restapi/operations"
)

var (
	_true = true
	_false = false
	apiResultOk = models.APIResult{Ok: &_true}
	apiResultError = models.APIResult{Ok: &_false}
)

//go:generate swagger generate server --target ../../chess2pic --name Chess2picAPI --spec ../api/chess2pic-api.yaml --principal interface{}

func configureFlags(api *operations.Chess2picAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.Chess2picAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.PostFenHandler = operations.PostFenHandlerFunc(func(params operations.PostFenParams) middleware.Responder {
		var from chess.PieceColor
		if *params.Body.FromWhite {
			from = chess.White
		} else {
			from = chess.Black
		}

		buf := &bytes.Buffer{}
		err := chess2pic.HandleFEN(strings.NewReader(*params.Body.Notation), buf, pic.DefaultCollection, from)
		if err != nil {
			resp := operations.NewPostFenBadRequest()
			resp.SetPayload(&models.ErrorResult{APIResult: apiResultError, Error: err.Error()})
			return resp
		} else {
			data := buf.Bytes()
			resp := operations.NewPostFenOK()
			resp.SetPayload(&models.ImageResult{APIResult: apiResultOk, Result: (*strfmt.Base64)(&data)})
			return resp
		}
	})
	api.PostPgnHandler = operations.PostPgnHandlerFunc(func(params operations.PostPgnParams) middleware.Responder {
		var from chess.PieceColor
		if *params.Body.FromWhite {
			from = chess.White
		} else {
			from = chess.Black
		}
		
		buf := &bytes.Buffer{}
		err := chess2pic.HandlePGN(strings.NewReader(*params.Body.Notation), buf, pic.DefaultCollection, from)
		if err != nil {
			resp := operations.NewPostPgnBadRequest()
			resp.SetPayload(&models.ErrorResult{APIResult: apiResultError, Error: err.Error()})
			return resp
		} else {
			data := buf.Bytes()
			resp := operations.NewPostPgnOK()
			resp.SetPayload(&models.ImageResult{APIResult: apiResultOk, Result: (*strfmt.Base64)(&data)})
			return resp
		}
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
