package mix

import (
	"io"
	"net/http"
)

//Write will set Headers and Write the response based on the Mixer
func Write(w http.ResponseWriter, m Mixer) error {
	for key, head := range m.Headers() {
		w.Header().Set(key, head)
	}

	_, err := io.Copy(w, m.Reader())
	return err
}

//Write will set Headers and Write the response based on the Mixer
func WriteStatus(w http.ResponseWriter, status int, m Mixer) error {
	for key, head := range m.Headers() {
		w.Header().Set(key, head)
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	_, err := io.Copy(w, m.Reader())
	return err
}
