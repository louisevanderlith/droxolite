package do

import "strconv"

type k struct {
	Name  string
	Clean bool
}

var serviceKeys map[k]string

func init() {
	serviceKeys = make(map[k]string)

	serviceKeys[k{"Router.API", false}] = "http://RouterAPI:8080/"
}

//GetServiceURL returns the correct URL for a service according to the caller's environment.
func GetServiceURL(instanceID, serviceName string, cleanURL bool) (string, error) {
	cacheService, ok := serviceKeys[k{serviceName, cleanURL}]

	if !ok {
		result := ""
		code, err := GET("", &result, instanceID, "Router.API", "discovery", instanceID, serviceName, strconv.FormatBool(cleanURL))

		if err != nil {
			scode := strconv.Itoa(code)
			return scode, err
		}

		serviceKeys[k{serviceName, cleanURL}] = result

		return result, nil
	}

	return cacheService, nil
}