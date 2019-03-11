package tools

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
)


func Dict(x interface{}) map[string]reflect.Value {
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)

	var methods map[string]reflect.Value
	methods = make(map[string]reflect.Value)
	for i := 0; i < v.NumMethod();i++ {
		methods[t.Method(i).Name] = v.Method(i)
	}
	return methods
}

func Dir(x interface{}) []string {
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)

	var methods []string
	for i := 0; i < v.NumMethod();i++ {
		methods = append(methods,t.Method(i).Name)
	}
	return methods
}

func GetAttr(x interface{},name string) reflect.Value {
	v := reflect.ValueOf(x)
	return v.MethodByName(name)
}

func Pop(x interface{}) (last interface{}, qq []interface{}){
	qq,_ = takeSliceArg(x)

	last = qq[len(qq)-1]
	qq = qq[0:len(qq)-1]
	return last, qq

}

func takeSliceArg(arg interface{}) (out []interface{},ok bool) {
	slice, success := takeArg(arg,reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0;i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func ReCompileGroup(rexp *regexp.Regexp,matchStr string) map[string]string {
	matchs := rexp.FindStringSubmatch(matchStr)
	if len(matchs) == 0 {
		return nil
	}

	groupNames := rexp.SubexpNames()
	result := make(map[string]string)
	for i,name := range groupNames {
		if i != 0 && name != "" {
			fmt.Println(name)
			fmt.Println(matchs)
			result[name] = matchs[i]
		}
	}
	return result
}

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}