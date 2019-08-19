package webserver

import (
	"net/http"

	boardModel "github.com/joostvdg/cmg/pkg/model"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
)

// GetMapLegend retrieves the map Legend, helps explain codes used within the data returned by the API
func GetMapLegend(c echo.Context) error {
	callback := c.QueryParam("callback")
	jsonp := c.QueryParam("jsonp")

	landscapes := make([]boardModel.Landscape, len(boardModel.Landscapes))
	i := 0
	for _, landscape := range boardModel.Landscapes {
		landscapes[i] = landscape
		i++
	}

	harbors := make([]boardModel.Harbor, len(boardModel.Harbors))
	j := 0
	for _, harbor := range boardModel.Harbors {
		harbors[j] = harbor
		j++
	}

	var content = model.MapLegend{
		Harbors:    harbors,
		Landscapes: landscapes,
	}

	if jsonp == "true" {
		return c.JSONP(http.StatusOK, callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}
