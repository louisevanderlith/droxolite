package sample

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type Stock struct {
}

func (req *Stock) Default(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Stock Home", false)

	return http.StatusOK, nil
}

type Parts struct {
}

func (req *Parts) Default(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Parts Default", false)

	return http.StatusOK, nil
}

func (req *Parts) Search(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Parts Search", false)

	return http.StatusOK, nil
}

func (req *Parts) View(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Parts View", false)

	return http.StatusOK, nil
}

func (req *Parts) Create(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Parts Create", false)

	return http.StatusOK, nil
}

type Services struct {
}

func (req *Services) Default(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Services Default", false)

	return http.StatusOK, nil
}

func (req *Services) Search(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Services Search", false)

	return http.StatusOK, nil
}

func (req *Services) View(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Services View", false)

	return http.StatusOK, nil
}

func (req *Services) Create(ctx context.Contexer) (int, interface{}) {
	//req.Setup("stock", "Services Create", false)

	return http.StatusOK, nil
}
