package components

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Info struct.
type Info struct {
	Name          string
	Dir           string
	GoPath        string
	TemplatePath  string
	AsyncJsPaths  []string
	DeferJsPaths  []string
	JsPaths       []string
	AsyncCssPaths []string
	CssPaths      []string
}

// Statics returns a slice of all static component paths.
func (info *Info) Statics() []string {
	list := make([]string, 0)
	list = append(list, info.AsyncJsPaths...)
	list = append(list, info.DeferJsPaths...)
	list = append(list, info.JsPaths...)
	list = append(list, info.AsyncCssPaths...)
	list = append(list, info.CssPaths...)
	return list
}

// List returns a list of components info.
func List(path string) ([]*Info, error) {
	list := make([]*Info, 0)
	if err := filepath.Walk(
		path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "component.go" {
				dir := filepath.Dir(p)
				name := strings.TrimPrefix(dir, path+"/")
				i, err := getInfo(name, dir)
				if err != nil {
					return err
				}
				list = append(list, i)
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return list, nil
}

func getInfo(name, dir string) (*Info, error) {
	info := &Info{
		Name:          name,
		Dir:           dir,
		GoPath:        filepath.Join(dir, "component.go"),
		TemplatePath:  filepath.Join(dir, "component.html"),
		AsyncJsPaths:  make([]string, 0),
		DeferJsPaths:  make([]string, 0),
		JsPaths:       make([]string, 0),
		AsyncCssPaths: make([]string, 0),
		CssPaths:      make([]string, 0),
	}
	goFound := false
	templateFound := false
	if err := filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i.IsDir() {
			return nil
		}
		if i.Name() == "component.go" {
			goFound = true
			return nil
		}
		if i.Name() == "component.html" {
			templateFound = true
			return nil
		}
		if strings.HasSuffix(i.Name(), ".async.js") {
			info.AsyncJsPaths = append(info.AsyncJsPaths, p)
			return nil
		}
		if strings.HasSuffix(i.Name(), ".defer.js") {
			info.DeferJsPaths = append(info.DeferJsPaths, p)
			return nil
		}
		if strings.HasSuffix(i.Name(), ".js") {
			info.JsPaths = append(info.JsPaths, p)
			return nil
		}
		if strings.HasSuffix(i.Name(), ".async.css") {
			info.AsyncCssPaths = append(info.AsyncCssPaths, p)
			return nil
		}
		if strings.HasSuffix(i.Name(), ".css") {
			info.CssPaths = append(info.CssPaths, p)
			return nil
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if !goFound {
		return nil, fmt.Errorf(
			"component %s does not contain component.go file in %s", name, dir)
	}
	if !templateFound {
		return nil, fmt.Errorf(
			"component %s does not contain component.html file in %s", name, dir)
	}
	return info, nil
}

func getPath() (string, error) {
	current, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	return filepath.Join(current, "components"), nil
}

func getHash() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type filesList []string

func (fl filesList) String() string {
	list := make([]string, len(fl))
	for i, s := range fl {
		list[i] = fmt.Sprintf(`"%s"`, s)
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(list, ","))
}
