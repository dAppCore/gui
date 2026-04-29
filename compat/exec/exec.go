package exec

import (
	"syscall"

	core "dappco.re/go"
)

type Cmd = core.Cmd

func Command(name string, arg ...string) *Cmd {
	return command(name, arg...)
}

func CommandContext(_ core.Context, name string, arg ...string) *Cmd {
	return command(name, arg...)
}

func LookPath(file string) (string, error) {
	if file == "" {
		return "", core.NewError("executable file not found")
	}
	if core.Contains(file, string(core.PathSeparator)) {
		if r := core.Stat(file); r.OK {
			return file, nil
		}
		return "", core.NewError("executable file not found")
	}
	for _, dir := range core.Split(core.Getenv("PATH"), string(core.PathListSeparator)) {
		if dir == "" {
			continue
		}
		candidate := core.PathJoin(dir, file)
		if r := core.Stat(candidate); r.OK {
			return candidate, nil
		}
	}
	return "", core.NewError("executable file not found")
}

func command(name string, arg ...string) *Cmd {
	path := name
	if resolved, err := LookPath(name); err == nil {
		path = resolved
	}
	return &Cmd{
		Path:        path,
		Args:        append([]string{name}, arg...),
		SysProcAttr: &syscall.SysProcAttr{},
	}
}
