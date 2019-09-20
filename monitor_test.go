package echoprometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
)

func TestMonitor(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewMetric()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	h(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}
