package command

import (
	"strconv"
)

type Arg struct {
	name  string
	usage string
	value string
}
type Flag struct {
	args map[string]*Arg
}

func (f *Flag) Var(name, defaultValue, usage string) {
	if f.args == nil {
		f.args = make(map[string]*Arg)
	}
	f.args[name] = &Arg{
		name:  name,
		usage: usage,
		value: defaultValue,
	}
}
func (f Flag) Get(name string) (string, bool) {
	if a, ok := f.args[name]; ok {
		return a.value, true
	} else {
		return "", false
	}
}
func (f Flag) GetInt(name string) (int, bool) {
	if s, ok := f.Get(name); !ok {
		return 0, ok
	} else {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, false
		}
		return n, true
	}
}
func (f Flag) Range(fun func(name, value, usage string)) {
	for _, v := range f.args {
		fun(v.name, v.value, v.usage)
	}
}
func (f *Flag) set(name, value string) {
	if a, ok := f.args[name]; ok {
		a.value = value
	} else {
		f.Var(name, value, "")
	}
}
