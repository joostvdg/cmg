package webserver

import (
	"fmt"
	"net/http"

	boardModel "github.com/joostvdg/cmg/pkg/model"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
)

// GetMapLegend retrieves the map Legend, helps explain codes used within the data returned by the API
func GetMapLegend(c echo.Context) error {
	callback := c.QueryParam("callback")
	jsonp := c.QueryParam("jsonp")

	harbors := make([]model.ResourceIdentity, 7)
	harbors[0] = model.ResourceIdentity{Name: "Grain", Id: fmt.Sprintf("%v", boardModel.Grain)}
	harbors[1] = model.ResourceIdentity{Name: "Wool", Id: fmt.Sprintf("%v", boardModel.Wool)}
	harbors[2] = model.ResourceIdentity{Name: "Ore", Id: fmt.Sprintf("%v", boardModel.Ore)}
	harbors[3] = model.ResourceIdentity{Name: "None", Id: fmt.Sprintf("%v", boardModel.None)}
	harbors[4] = model.ResourceIdentity{Name: "All", Id: fmt.Sprintf("%v", boardModel.All)}
	harbors[5] = model.ResourceIdentity{Name: "Brick", Id: fmt.Sprintf("%v", boardModel.Brick)}
	harbors[6] = model.ResourceIdentity{Name: "Lumber", Id: fmt.Sprintf("%v", boardModel.Lumber)}

	landscapes := make([]model.ResourceIdentity, 6)
	landscapes[0] = model.ResourceIdentity{Name: "Desert", Id: fmt.Sprintf("%v", boardModel.Desert)}
	landscapes[1] = model.ResourceIdentity{Name: "Forest", Id: fmt.Sprintf("%v", boardModel.Forest)}
	landscapes[2] = model.ResourceIdentity{Name: "Pasture", Id: fmt.Sprintf("%v", boardModel.Pasture)}
	landscapes[3] = model.ResourceIdentity{Name: "Field", Id: fmt.Sprintf("%v", boardModel.Field)}
	landscapes[4] = model.ResourceIdentity{Name: "River", Id: fmt.Sprintf("%v", boardModel.River)}
	landscapes[5] = model.ResourceIdentity{Name: "Mountain", Id: fmt.Sprintf("%v", boardModel.Mountain)}

	var content = model.MapLegend{
		Harbors:    harbors,
		Landscapes: landscapes,
	}

	if jsonp == "true" {
		return c.JSONP(http.StatusOK, callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}
