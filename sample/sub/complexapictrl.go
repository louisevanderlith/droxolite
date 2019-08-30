package sub

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type ComplexAPICtrl struct {
}

func (c *ComplexAPICtrl) Get(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, "This is complex!"
}
