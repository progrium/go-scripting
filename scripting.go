package scripting

import (
	"sync"
	"strings"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/progrium/go-extensions"
)

type RuntimeEngine interface {
	FileExtension() string
	CallModule(name, function string, args []interface{}) (interface{}, error)
	InitModule(name, source string, globals map[string]interface{}) error
	UpdateGlobals(globals map[string]interface{})
}

var Runtimes = extensions.ExtensionPoint(new(RuntimeEngine))

var scripting = struct {
	sync.Mutex
	globals map[string]interface{}
	modules map[string]RuntimeEngine
} {
	globals: make(map[string]interface{}),
	modules: make(map[string]RuntimeEngine),
}

func UpdateGlobals(globals map[string]interface{}) {
	scripting.Lock()
	defer scripting.Unlock()
	for k, v := range globals {
		scripting.globals[k] = v
	}
	for _, r := range Runtimes.All() {
		r.(RuntimeEngine).UpdateGlobals(scripting.globals)
	}
}

func GetGlobal(name string) interface{} {
	scripting.Lock()
	defer scripting.Unlock()
	return scripting.globals[name]
}

func LoadModule(name, source string, runtime string) error {
	scripting.Lock()
	defer scripting.Unlock()
	r := Runtimes.Get(runtime).(RuntimeEngine)
	err := r.InitModule(name, source, scripting.globals)
	if err != nil {
		return err
	}
	scripting.modules[name] = r
	return nil
}

func findRuntimeForFile(path string) string {
	for name, runtime := range Runtimes.All() {
		fileExt := runtime.(RuntimeEngine).FileExtension()
		if fileExt != "" && strings.HasSuffix(path, fileExt) {
			return name
		}
	}
	return ""
}

func LoadModuleFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	runtime := findRuntimeForFile(path)
	if runtime == "" {
		return errors.New("scripting: no runtime found to handle: " + path)
	}
	name := strings.Split(filepath.Base(path), ".")[0]
	return LoadModule(name, string(data), runtime)
}

func LoadModulesFromPath(path string) error {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		filepath := path + "/" + entry.Name()
		runtime := findRuntimeForFile(filepath)
		if runtime != "" {
			err = LoadModuleFile(filepath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Call(module, function string, args []interface{}) (interface{}, error) {
	scripting.Lock()
	runtime, ok := scripting.modules[module]
	scripting.Unlock()
	if !ok {
		return nil, errors.New("scripting: no such module loaded: " + module)
	}
	return runtime.(RuntimeEngine).CallModule(module, function, args)
}
