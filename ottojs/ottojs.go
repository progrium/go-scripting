package ottojs

import (
	"reflect"
	"sync"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/progrium/go-scripting"
)

func Register() {
	scripting.Runtimes.Register("js", &RuntimeEngine{
		modules: make(map[string]*otto.Otto),
	})
}

type RuntimeEngine struct {
	sync.Mutex
	modules map[string]*otto.Otto
}

func (r *RuntimeEngine) FileExtension() string {
	return ".js"
}

func (r *RuntimeEngine) InitModule(name, source string, globals map[string]interface{}) error {
	r.Lock()
	defer r.Unlock()
	r.modules[name] = otto.New()
	r.setModuleGlobals(name, globals)
	r.modules[name].Run(source)
	return nil
}

func (r *RuntimeEngine) CallModule(name, function string, args []interface{}) (interface{}, error) {
	r.Lock()
	context := r.modules[name]
	r.Unlock()
	value, err := context.Call(function, nil, args...)
	if err != nil {
		return nil, err
	}
	exported, _ := value.Export()
	return exported, nil
}

func (r *RuntimeEngine) setModuleGlobals(name string, globals map[string]interface{}) {
	context := r.modules[name]
	for k, v := range globals {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			setValueAtPath(context, k, funcToOtto(context, reflect.ValueOf(v)))
		} else {
			setValueAtPath(context, k, v)
		}
	}
}

func (r *RuntimeEngine) UpdateGlobals(globals map[string]interface{}) {
	r.Lock()
	defer r.Unlock()
	for module := range r.modules {
		r.setModuleGlobals(module, globals)
	}
}

func setValueAtPath(context *otto.Otto, path string, value interface{}) {
	parts := strings.Split(path, ".")
	parentCount := len(parts) - 1
	if parentCount > 0 {
		parentPath := strings.Join(parts[0:parentCount], ".")
		parent, err := context.Object("(" + parentPath + ")")
		if err != nil {
			emptyObject, _ := context.Object(`({})`)
			setValueAtPath(context, parentPath, emptyObject)
		}
		parent, _ = context.Object("(" + parentPath + ")")
		parent.Set(parts[parentCount], value)
	} else {
		context.Set(path, value)
	}
}

func funcToOtto(context *otto.Otto, fn reflect.Value) interface{} {
	return func(call otto.FunctionCall) otto.Value {
		convertedArgs := make([]reflect.Value, 0)
		for _, v := range call.ArgumentList {
			exported, _ := v.Export()
			convertedArgs = append(convertedArgs, reflect.ValueOf(exported))
		}
		ret := fn.Call(convertedArgs)
		if len(ret) > 0 {
			val, _ := context.ToValue(ret[0].Interface())
			return val
		} else {
			return otto.UndefinedValue()
		}
	}
}