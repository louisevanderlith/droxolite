package models

import (
	"errors"
	"strings"

	"github.com/louisevanderlith/husk"
)

// Client contains the values required for dynamic client registration
// https://tools.ietf.org/html/rfc7591
type Client struct {
	Name               string  `hsk:"size(50)"` //clientname
	Description        string  `hsk:"size(256)"`
	RedirectURI        string  `hsk:"null"` //must be https only if auth_code flow
	EndpointAuthMethod string  //none (public), client_secret_post, client_secret_basic
	GrantType          string  //authorization_code (public), implicit (not supported), password, client_credentials, refresh_token
	ResponseType       string  //code, token, (none)
	ClientURI          string  //used for display only
	LogoURI            string  `hsk:"null"` //logo,
	Scopes             []Scope //Update to constraint relationship
	Folio              string  `hsk:"size(50);null"` //(optional) name of the folio profile this client uses.
	Secret             string  `json:"-"`
}

// NewPublicClient returns a client which follows the 'auth_code' flow, usually user facing applications
func NewPublicClient(name, descr, uri, redirect string) Client {
	return Client{
		Name:               name,
		Description:        descr,
		RedirectURI:        redirect,
		EndpointAuthMethod: "none",
		GrantType:          "authorization_code",
		ResponseType:       "code",
		ClientURI:          uri,
	}
}

// NewPrivateClient returns a client which follows the 'resoure owner/password' flow, used by web services
func NewPrivateClient(name, descr, uri string) Client {
	return Client{
		Name:        name,
		Description: descr,
		//No Redirect Required
		EndpointAuthMethod: "client_secret_basic",
		GrantType:          "password",
		ResponseType:       "token",
		ClientURI:          uri,
	}
}

func (c Client) Valid() (bool, error) {

	if len(c.Name) == 0 {
		return false, errors.New("client name required")
	}

	if c.GrantType == "implicit" {
		return false, errors.New("implicit is not supported for security reasons")
	}

	if len(c.EndpointAuthMethod) == 0 {
		return false, errors.New("endpoint auth method required")
	}

	if c.EndpointAuthMethod == "none" {
		if !strings.HasPrefix(c.RedirectURI, "https://") {
			return false, errors.New("invalid redirect URI")
		}

		hasGrant := c.GrantType == "authorization_code"

		if !hasGrant {
			return false, errors.New("invalid grant type")
		} else if c.ResponseType != "code" {
			return false, errors.New("invalid response type")
		}
	}

	if len(c.Scopes) == 0 {
		return false, errors.New("atleast one scope must be defined")
	}

	return husk.ValidateStruct(&c)
}

/*

//NewService returns a new instance of a Services' information
//publicKey refers to the location of the public key file (.pub)
func NewService(name, secret, profile, publicKey, host string, port int, serviceType servicetype.Enum) *Service {
	result := &Service{
		Name:           fmt.Sprintf("%s.%s", name, serviceType),
		Type:           serviceType,
		PublicKey:      publicKey,
		AllowedCallers: make(map[servicetype.Enum]struct{}),
		Port:           port,
		Profile:        profile,
		Host:           host,
		Secret:         secret,
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

*/
