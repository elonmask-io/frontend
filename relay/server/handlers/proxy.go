package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
}

func (p *ProxyHandler) Proxy(c echo.Context) error {
	address := c.Get("address").(string)

	uri, err := url.Parse(address + c.Path())
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		if errors.Is(context.Canceled, err) {
			httpError := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("client closed connection: %w", err))
			httpError.Internal = err
			c.Set("_error", httpError)
		} else if err != nil {
			httpError := echo.NewHTTPError(http.StatusBadGateway, fmt.Errorf("enclave is unreachable: %w", err))
			httpError.Internal = err
			c.Set("_error", httpError)
		}
	}

	proxy.ServeHTTP(c.Response(), c.Request())

	return nil
}