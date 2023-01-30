package comment

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
func (f *Flag) set(name, value string) {
	if a, ok := f.args[name]; ok {
		a.value = value
	}
}
