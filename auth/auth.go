// Package auth provides functionality to fetch authentication information for
// various Docker registries. It provides special handling for Docker Hub.
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

const ENV_NAME = "INVOLUCRO_AUTH"

type authenticationInfo struct {
	docker.AuthConfiguration
}

func (a *authenticationInfo) UnmarshalString(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.User == nil {
		return fmt.Errorf("invalid authentication URI")
	}
	a.Username = u.User.Username()
	a.Password, _ = u.User.Password()
	a.ServerAddress = u.Host + u.Path
	a.Email = u.Query().Get("email")
	return nil
}

func (a *authenticationInfo) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}
	return a.UnmarshalString(s)
}

func getAllFrom(r io.Reader) ([]authenticationInfo, error) {
	res := struct {
		AuthenticationInfos []authenticationInfo `json:"auths"`
	}{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&res); err != nil {
		return []authenticationInfo{}, err
	}
	return res.AuthenticationInfos, nil
}

// ForServer reads the configuration file and returns a AuthConfiguration, when
// there is one for the given server.
//
// If there is one, the second return value will be true, if there is none, it
// is false.
//
// However, if an error other than file not found occurs, this error will be
// returned and the value of the other values is undefined.
func ForServer(server string) (docker.AuthConfiguration, bool, error) {
	filename := path.Join(userHomeDir(), ".involucro")

	file, err := os.Open(filename)
	if err != nil {
		return docker.AuthConfiguration{}, false, nil
	}
	defer file.Close()
	return forServerInFile(server, file)
}

func forServerInFile(server string, file io.Reader) (docker.AuthConfiguration, bool, error) {
	env := os.Getenv(ENV_NAME)
	if env != "" {
		as := strings.Split(env, " ")
		ai := authenticationInfo{}

		for _, a := range as {
			err := ai.UnmarshalString(a)
			if err != nil {
				return docker.AuthConfiguration{}, false, err
			}
			if ai.ServerAddress == server {
				if server == "index.docker.io/v1/" {
					ai.ServerAddress = "https://index.docker.io/v1/"
				}
				return ai.AuthConfiguration, true, nil
			}
		}
	}

	if server == "" {
		server = "index.docker.io/v1/"
	}
	allEntries, err := getAllFrom(file)
	if err != nil {
		return docker.AuthConfiguration{}, false, err
	}
	for _, el := range allEntries {
		if el.ServerAddress == server {
			if server == "index.docker.io/v1/" {
				el.ServerAddress = "https://index.docker.io/v1/"
			}
			return el.AuthConfiguration, true, nil
		}
	}

	return docker.AuthConfiguration{}, false, nil
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
