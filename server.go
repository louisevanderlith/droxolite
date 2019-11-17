package droxolite

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/louisevanderlith/droxolite/resins"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	readTimeout  = time.Second * 15
	writeTimeout = time.Second * 15
)

//Boot starts the Epoxy Objects to serve a configured service.
func Boot(e resins.Epoxi, port int) error {
	srvr := newServer(port)
	srvr.Handler = e.Router()

	return srvr.ListenAndServe()
}

//Boot starts the Epoxy Objects to securely serve a configured service
func BootSecure(e resins.Epoxi, pubKeyPath, privKeyPath string, port, fromPort int) error {
	publicKeyPem := readBlocks(pubKeyPath)
	privateKeyPem := readBlocks(privKeyPath)
	cert, err := tls.X509KeyPair(publicKeyPem, privateKeyPem)

	if err != nil {
		return err
	}

	srvr := newServer(port)
	srvr.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	srvr.Handler = e.Router()

	err = srvr.ListenAndServeTLS("", "")

	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", fromPort), http.HandlerFunc(redirectTLS))
}

//GenerateSecret returns a Base-16 UUID
func GenerateSecret() string {
	id, err := uuid.NewV4()

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", id)
}

func newServer(port int) *http.Server {
	return &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         fmt.Sprintf(":%v", port),
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	moveURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, moveURL, http.StatusPermanentRedirect)
}

func readBlocks(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
