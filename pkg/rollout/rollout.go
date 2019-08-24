package rollout

import (
	roxServer "github.com/rollout/rox-go/server"
)

type Container struct {
	EnableTutorial         roxServer.RoxFlag
	EnableHarborValidation roxServer.RoxFlag
}

var Rox *roxServer.Rox
var RoxContainer = &Container{
	EnableTutorial:         roxServer.NewRoxFlag(false),
	EnableHarborValidation: roxServer.NewRoxFlag(false),
}
