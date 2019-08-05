package sample

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type Stock struct {
	xontrols.UICtrl
}

func (req *Stock) Default() {
	req.Setup("stock", "Stock Home", false)

	req.Serve(http.StatusOK, nil, nil)
}

type Parts struct {
	xontrols.UICtrl
}

func (req *Parts) Default() {
	req.Setup("stock", "Parts Default", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) Search() {
	req.Setup("stock", "Parts Search", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) View() {
	req.Setup("stock", "Parts View", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) Create() {
	req.Setup("stock", "Parts Create", false)

	req.Serve(http.StatusOK, nil, nil)
}

type Services struct {
	xontrols.UICtrl
}

func (req *Services) Default() {
	req.Setup("stock", "Services Default", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) Search() {
	req.Setup("stock", "Services Search", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) View() {
	req.Setup("stock", "Services View", false)

	req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) Create() {
	req.Setup("stock", "Services Create", false)

	req.Serve(http.StatusOK, nil, nil)
}
