package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetMapLegend(t *testing.T) {
	targetPath := "/api/legend"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, GetMapLegend(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var mapLegend model.MapLegend

	expectedNumberOfHarbors := 7
	expectedNumberOfResources := 6

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &mapLegend)) {
		assert.Equal(t, len(mapLegend.Harbors), expectedNumberOfHarbors)
		assert.Equal(t, len(mapLegend.Landscapes), expectedNumberOfResources)
	}
}

func TestGetMapLegendSupportsJsonP(t *testing.T) {
	targetPath := "/api/legend?jsonp=true"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, GetMapLegend(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var mapLegend model.MapLegend

	fmt.Printf("Body: %v\n", rec.Body.String())

	assert.True(t, strings.HasPrefix(rec.Body.String(), "("))
	assert.True(t, strings.HasSuffix(rec.Body.String(), ");"))
	strippedBody := strings.TrimPrefix(rec.Body.String(), "(")
	strippedBody = strings.TrimSuffix(strippedBody, ");")

	assert.NoError(t, json.Unmarshal([]byte(strippedBody), &mapLegend))
}

func TestGetMapLegendSupportsJsonPAndCallback(t *testing.T) {
	callback := "MyFunction"
	targetPath := fmt.Sprintf("/api/legend?jsonp=true&callback=%v", callback)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, GetMapLegend(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()

	prefix := fmt.Sprintf("%v(", callback)
	assert.True(t, strings.HasPrefix(rec.Body.String(), prefix))
	assert.True(t, strings.HasSuffix(rec.Body.String(), ");"))
	strippedBody := strings.TrimPrefix(rec.Body.String(), prefix)
	strippedBody = strings.TrimSuffix(strippedBody, ");")
}
