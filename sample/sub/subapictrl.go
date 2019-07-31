package sub

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type SubAPICtrl struct {
	xontrols.APICtrl
}

func (c *SubAPICtrl) Get() {
	c.Serve(http.StatusOK, nil, "I am a sub controller")
}
