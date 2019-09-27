package bodies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/louisevanderlith/droxolite/servicetype"
)

//Service identifies the Registering APP or API
type Service struct {
	ID             string
	Name           string
	URL            string
	PublicURL      string
	Host           string
	Version        int
	AllowedCallers map[servicetype.Enum]struct{}
	Type           servicetype.Enum
	PublicKey      string
	Port           int
	Profile        string
}

//NewService returns a new instance of a Services' information
//publicKey refers to the location of the public key file (.pub)
func NewService(name, profile, publicKey, host string, port int, serviceType servicetype.Enum) *Service {
	result := &Service{
		Name:           fmt.Sprintf("%s.%s", name, serviceType),
		Type:           serviceType,
		PublicKey:      publicKey,
		AllowedCallers: make(map[servicetype.Enum]struct{}),
		Port:           port,
		Profile:        profile,
		Host:           host,
	}

	return result
}

// Register is used to register an application with the router service
func (s *Service) Register(routerUrl string) error {
	err := s.setURL()

	if err != nil {
		return err
	}

	resp, err := s.sendRegistration(routerUrl)

	if err != nil {
		return err
	}

	if len(resp.Reason) > 0 {
		return resp
	}

	s.ID = resp.Data.(string)

	return nil
}

func (s *Service) sendRegistration(routerUrl string) (*RESTResult, error) {
	bits, err := json.Marshal(s)

	if err != nil {
		return nil, err
	}

	disco := fmt.Sprintf("%sdiscovery", routerUrl)
	resp, err := http.Post(disco, "application/json", bytes.NewBuffer(bits))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	data, err := MarshalToResult(contents, "")

	return data, err
}

func (s *Service) setURL() error {
	url, err := getNetworkIP(s.Name, strconv.Itoa(s.Port))

	if err != nil {
		return err
	}

	s.URL = url
	s.PublicURL = makeURL(s.Host, strconv.Itoa(s.Port))

	return nil
}

func getNetworkIP(name, port string) (string, error) {
	uniqueName := strings.Replace(name, ".", "", 1)

	return makeURL(uniqueName, port), nil
}

func makeURL(domain, port string) string {
	schema := "http"

	return fmt.Sprintf("%s://%s:%s/", schema, domain, port)
}
