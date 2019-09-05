package clients

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

func Index(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, "You're Home!"
}
