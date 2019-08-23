package sub

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type SubAPICtrl struct {
}

func (c *SubAPICtrl) Get(ctx context.Contexer) (int, interface{}) {
	return http.StatusOK, "I am a sub controller"
}
