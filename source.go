package i18n // import "github.com/novakit/i18n"

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// Source i18n source
type Source struct {
	loaded bool
	values map[string]string
	fs     http.FileSystem
	l      *sync.RWMutex
}

// NewSource creaet a new i18n source
func NewSource(fs http.FileSystem) *Source {
	return &Source{
		values: map[string]string{},
		fs:     fs,
		l:      &sync.RWMutex{},
	}
}

// Flatten flatten a nested map[interface{}]interface{} into a flatten map[string]string with dot combined key
func Flatten(pfx string, in map[interface{}]interface{}, out map[string]string) {
	for key, val := range in {
		keyStr := fmt.Sprintf("%v", key)
		if valMap, ok := val.(map[interface{}]interface{}); ok {
			Flatten(pfx+keyStr+".", valMap, out)
		} else if valStr, ok := val.(string); ok {
			out[pfx+keyStr] = valStr
		} else {
			out[pfx+keyStr] = fmt.Sprintf("%v", val)
		}
	}
}

func (s *Source) loadYAML(pfx string, buf []byte) (err error) {
	d := map[interface{}]interface{}{}
	if err = yaml.Unmarshal(buf, &d); err != nil {
		return
	}
	Flatten(pfx, d, s.values)
	return
}

func (s *Source) load() (err error) {
	var f http.File
	if f, err = s.fs.Open("."); err != nil {
		return
	}
	defer f.Close()
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return
	}
	if !fi.IsDir() {
		err = fmt.Errorf("i18n: is not a directory")
		return
	}
	var fis []os.FileInfo
	if fis, err = f.Readdir(-1); err != nil {
		return
	}
	for _, fi1 := range fis {
		fn1 := fi1.Name()
		ext := path.Ext(fn1)
		if ext != ".yml" && ext != ".yaml" {
			continue
		}
		// extract locale name
		n := fn1[:len(fn1)-len(ext)]
		if strings.HasPrefix(n, "/") {
			n = n[1:]
		}
		// merge en-US.aaa.yml with en-US.bbb.yml
		if ns := strings.Split(n, "."); len(ns) > 1 {
			n = ns[0]
		}
		// open file
		var f1 http.File
		if f1, err = s.fs.Open(fn1); err != nil {
			break
		}
		var buf []byte
		buf, err = ioutil.ReadAll(f1)
		if err == nil {
			err = s.loadYAML(n+".", buf)
		}
		f1.Close()
		if err != nil {
			break
		}
	}
	return
}

// LoadIfNeeded load source if not already loaded
func (s *Source) LoadIfNeeded() {
	if s.loaded {
		return
	}
	s.l.Lock()
	defer s.l.Unlock()
	s.load()
	s.loaded = true
}

// Get get a value by key
func (s *Source) Get(key string) string {
	s.LoadIfNeeded()
	s.l.RLock()
	defer s.l.RUnlock()
	return s.values[key]
}

// Reload reload all values
func (s *Source) Reload() {
	s.l.Lock()
	defer s.l.Unlock()
	s.loaded = false
	s.values = map[string]string{}
}
