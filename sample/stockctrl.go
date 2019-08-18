package sample

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type Stock struct {
	xontrols.UICtrl
}

func (req *Stock) Default() error {
	req.Setup("stock", "Stock Home", false)

	return req.Serve(http.StatusOK, nil, nil)
}

type Parts struct {
	xontrols.UICtrl
}

func (req *Parts) Default() error {
	req.Setup("stock", "Parts Default", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) Search() error {
	req.Setup("stock", "Parts Search", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) View() error {
	req.Setup("stock", "Parts View", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Parts) Create() error {
	req.Setup("stock", "Parts Create", false)

	return req.Serve(http.StatusOK, nil, nil)
}

type Services struct {
	xontrols.UICtrl
}

func (req *Services) Default() error {
	req.Setup("stock", "Services Default", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) Search() error {
	req.Setup("stock", "Services Search", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) View() error {
	req.Setup("stock", "Services View", false)

	return req.Serve(http.StatusOK, nil, nil)
}

func (req *Services) Create() error {
	req.Setup("stock", "Services Create", false)

	return req.Serve(http.StatusOK, nil, nil)
}
