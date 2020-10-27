package mod

import (
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	rx = regexp.MustCompile("^module [\"]{0,1}(.+)[\"]{0,1}$")
)

// FromGoMod returns package name from go.mod file.
func FromGoMod(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, row := range strings.Split(string(b)+"\n", "\n") {
		if v := rx.FindAllStringSubmatch(row, 1); len(v) == 1 {
			if len(v[0]) == 2 {
				return v[0][1], nil
			}
			continue
		}
	}
	return "", nil
}
