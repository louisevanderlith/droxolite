package clients

import (
	"fmt"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type Nomad struct {
}

func (x *Nomad) AcceptsQuery() map[string]string {
	result := make(map[string]string)
	result["name"] = "{name}"
	return result
}

func (x *Nomad) Get(ctx context.Requester) (int, interface{}) {
	param := ctx.FindQueryParam("name")
	result := fmt.Sprintf("Nomad got %s", param)
	return http.StatusOK, result
}
