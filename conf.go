package droxolite

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

//Conf is used to load the data from '/conf/conf.json'
type Conf struct {
	Appname   string
	HTTPPort  int
	HTTPSPort int
	Profile   string
	UnitValue int
}

//LoadSettings returns the data contained in the 'domains.json' config file.
func LoadConfig() (*Conf, error) {
	dbConfPath := FindFilePath("conf.json", "conf")
	content, err := ioutil.ReadFile(dbConfPath)

	if err != nil {
		return nil, err
	}

	settings := &Conf{}
	err = json.Unmarshal(content, settings)

	if err != nil {
		return nil, err
	}

	return settings, nil
}

//Returns the filepath within the current working directory.
func FindFilePath(fileName, targetFolder string) string {
	var result string
	wp := getWorkingPath() + "/" + targetFolder

	result = filepath.Join(wp, filepath.FromSlash(path.Clean("/"+fileName)))

	return result
}

func getWorkingPath() string {
	ex, err := os.Getwd()

	if err != nil {
		log.Print("getWorkingPath: ", err)
	}

	return ex
}
