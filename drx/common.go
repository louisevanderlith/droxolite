package drx

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

//JSONBody will bind JSON body to container
func JSONBody(r *http.Request, container interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(container)
}

//FindParam will return named path parameter value
func FindParam(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}

//FindQueryParam returns the First requested querystring parameter
func FindQueryParam(r *http.Request, name string) string {
	results, ok := r.URL.Query()[name]

	if !ok {
		return ""
	}

	return results[0]
}

//GetPageData will attempt to parse the request Url for page sizes.
func GetPageData(r *http.Request) (page, pageSize int) {
	pageData := FindParam(r, "pagesize")
	return ParsePageData(pageData)
}

//turns /B5 into page 2. size 5
func ParsePageData(pageData string) (int, int) {
	defaultPage := 1
	defaultSize := 10

	if len(pageData) < 2 {
		return defaultPage, defaultSize
	}

	pChar := []rune(pageData[:1])

	if len(pChar) != 1 {
		return defaultPage, defaultSize
	}

	page := int(pChar[0]) % 32
	pageSize, err := strconv.Atoi(pageData[1:])

	if err != nil {
		return defaultPage, defaultSize
	}

	return page, pageSize
}

func GravatarHash(email string) string {
	if len(email) == 0 {
		return ""
	}

	gravatar := md5.Sum([]byte(strings.ToLower(strings.Trim(email, " "))))

	return fmt.Sprintf("%x", gravatar)
}
