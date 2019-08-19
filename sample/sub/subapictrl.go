package sub

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type SubAPICtrl struct {
	xontrols.APICtrl
}

func (c *SubAPICtrl) Get() error {
	return c.Serve(http.StatusOK, nil, "I am a sub controller")
}
