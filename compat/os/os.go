package os

import (
	"syscall"

	c "dappco.re/go"
)

type File = c.OSFile
type FileInfo = c.FsFileInfo
type FileMode = c.FileMode

const ModeSymlink = c.ModeSymlink

var (
	Args                    = c.Args()
	Stderr         c.Writer = c.Stderr()
	ErrNotExist             = syscall.ENOENT
	ErrProcessDone          = syscall.ECHILD
)

func Chdir(dir string) error {
	r := c.Chdir(dir)
	return resultError(r)
}

func Environ() []string {
	return c.Environ()
}

func Exit(code int) {
	c.Exit(code)
}

func Getenv(key string) string {
	return c.Getenv(key)
}

func Getpid() int {
	return c.Getpid()
}

func Hostname() (string, error) {
	r := c.Hostname()
	if !r.OK {
		return "", resultError(r)
	}
	return r.Value.(string), nil
}

func IsNotExist(err error) bool {
	return c.IsNotExist(err)
}

func Lstat(name string) (FileInfo, error) {
	r := c.Lstat(name)
	if !r.OK {
		return nil, resultError(r)
	}
	return r.Value.(FileInfo), nil
}

func LookupEnv(key string) (string, bool) {
	return c.LookupEnv(key)
}

func Mkdir(name string, perm FileMode) error {
	r := c.Mkdir(name, perm)
	return resultError(r)
}

func MkdirAll(path string, perm FileMode) error {
	r := c.MkdirAll(path, perm)
	return resultError(r)
}

func MkdirTemp(dir, pattern string) (string, error) {
	r := c.MkdirTemp(dir, pattern)
	if !r.OK {
		return "", resultError(r)
	}
	return r.Value.(string), nil
}

func Open(name string) (*File, error) {
	r := c.Open(name)
	if !r.OK {
		return nil, resultError(r)
	}
	return r.Value.(*File), nil
}

func ReadFile(name string) ([]byte, error) {
	r := c.ReadFile(name)
	if !r.OK {
		return nil, resultError(r)
	}
	return r.Value.([]byte), nil
}

func RemoveAll(path string) error {
	r := c.RemoveAll(path)
	return resultError(r)
}

func Setenv(key, value string) error {
	r := c.Setenv(key, value)
	return resultError(r)
}

func Stat(name string) (FileInfo, error) {
	r := c.Stat(name)
	if !r.OK {
		return nil, resultError(r)
	}
	return r.Value.(FileInfo), nil
}

func Symlink(oldname, newname string) error {
	return syscall.Symlink(oldname, newname)
}

func TempDir() string {
	return c.TempDir()
}

func Unsetenv(key string) error {
	r := c.Unsetenv(key)
	return resultError(r)
}

func UserCacheDir() (string, error) {
	r := c.UserCacheDir()
	if !r.OK {
		return "", resultError(r)
	}
	return r.Value.(string), nil
}

func UserConfigDir() (string, error) {
	r := c.UserConfigDir()
	if !r.OK {
		return "", resultError(r)
	}
	return r.Value.(string), nil
}

func UserHomeDir() (string, error) {
	r := c.UserHomeDir()
	if !r.OK {
		return "", resultError(r)
	}
	return r.Value.(string), nil
}

func WriteFile(name string, data []byte, perm FileMode) error {
	r := c.WriteFile(name, data, perm)
	return resultError(r)
}

func resultError(r c.Result) error {
	if r.OK {
		return nil
	}
	if err, ok := r.Value.(error); ok {
		return err
	}
	return c.NewError(r.Error())
}
