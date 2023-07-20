package http

import (
	"fmt"
	"net/url"
)

func BuildPath(base string, pathElements ...string) string {
	path, err := url.JoinPath(base, pathElements...)

	if err != nil {
		panic(fmt.Sprintf("Error while building path: %s", err))
	} else {
		return path
	}
}
