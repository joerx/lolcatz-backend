package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joerx/lolcatz-backend/http/middleware"
)

const corsOrigin = "example.com"

func Test_httpHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/"), nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected statuscode %d, but got %d", e, a)
	}

	// Content-type
	if e, a := "application/json", w.Result().Header.Get("Content-type"); e != a {
		t.Errorf("Expected content type to be '%s' but got '%s'", e, a)
	}
}

func Test_corsHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, fmt.Sprintf("/"), nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected statuscode %d, but got %d", e, a)
	}

	// Access-Control-Allow-Origin
	if e, a := corsOrigin, w.Result().Header.Get(middleware.HeaderAccessControlAllowOrigin); e != a {
		t.Errorf("Expected CORS origin to be '%s' but got '%s'", e, a)
	}

	// Access-Control-Allow-Methods
	expectedCORSMethods := []string{http.MethodPut, http.MethodDelete, http.MethodOptions}
	a := w.Result().Header.Get(middleware.HeaderAccessControlAllowMethods)
	for _, e := range expectedCORSMethods {
		if !strings.Contains(a, e) {
			t.Errorf("Expected Access-Control-Allowed-Methods to contain '%s' but only got '%s'", e, a)
		}
	}

	// Access-Control-Allow-Headers
	allowedHeaders := []string{"Accept", "Content-Type", "Content-Length", "Authorization"}
	a = w.Result().Header.Get(middleware.HeaderAccessControlAllowHeaders)
	for _, e := range allowedHeaders {
		if !strings.Contains(a, e) {
			t.Errorf("Expected Access-Control-Allow-Headers to contain '%s' but only got '%s'", e, a)
		}
	}
}
