package sub

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type ComplexAPICtrl struct {
	xontrols.APICtrl
}

func (c *ComplexAPICtrl) Get() {
	c.Serve(http.StatusOK, nil, "This is complex!")
}
