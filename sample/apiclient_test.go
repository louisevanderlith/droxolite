package sample

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisevanderlith/droxolite"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/servicetype"
)

func TestMain_API(t *testing.T) {
	srvc := droxolite.NewService("Test.API", "/certs/none.pem", 8090, servicetype.API)
	srvc.ID = "Tester1"

	poxy := droxolite.NewEpoxy(srvc)
	apiRoutes(poxy)

	req, err := http.NewRequest("GET", "/fake/", nil)

	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handle := poxy.GetRouter()
	handle.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "48"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func apiRoutes(poxy *droxolite.Epoxy) {
	fakeCtrl := &FakeAPICtrl{}

	fkgroup := droxolite.NewRouteGroup("Fake", fakeCtrl)
	fkgroup.AddRoute("/", "GET", roletype.Admin, fakeCtrl.Get)
	fkgroup.AddRoute("/{id:[0-9]+}", "GET", roletype.Admin, fakeCtrl.GetId)
	poxy.AddGroup(fkgroup)
}

/*
keyPath := os.Getenv("KEYPATH")
	pubName := os.Getenv("PUBLICKEY")
	pubPath := path.Join(keyPath, pubName)

	conf, err := domains.LoadConfig()

	if err != nil {
		panic(err)
	}

	name := conf.AppName

	srv := mango.NewService(name, pubPath, enums.APP)

	port := conf.HTTPPort
	err := srv.Register(port)

	if err != nil {
		log.Print("Register: ", err)
	} else {
		err = mango.UpdateTheme(srv.ID)

		if err != nil {
			panic(err)
		}

		routers.Setup(srv)

		//beego.SetStaticPath("/dist", "dist")
		//beego.Run()
	}
*/
