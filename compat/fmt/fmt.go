package fmt

import core "dappco.re/go"

func Sprint(args ...any) string {
	return core.Sprint(args...)
}

func Sprintf(format string, args ...any) string {
	return core.Sprintf(format, args...)
}

func Sprintln(args ...any) string {
	return core.Sprintln(args...)
}

func Errorf(format string, args ...any) error {
	return core.Errorf(format, args...)
}

func Println(args ...any) {
	core.Println(args...)
}

func Printf(format string, args ...any) {
	core.Print(core.Stdout(), format, args...)
}

func Fprintln(w core.Writer, args ...any) (int, error) {
	return w.Write([]byte(core.Sprintln(args...)))
}
