package filepath

import core "dappco.re/go"

const Separator = core.PathSeparator

var SkipDir = core.PathSkipDir
var SkipAll = core.PathSkipAll

type WalkFunc = core.PathWalkFunc

func Abs(path string) (string, error) {
	r := core.PathAbs(path)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return "", err
		}
		return "", core.NewError(r.Error())
	}
	return r.Value.(string), nil
}

func Base(path string) string {
	return core.PathBase(path)
}

func Clean(path string) string {
	return core.CleanPath(path, string(core.PathSeparator))
}

func Dir(path string) string {
	return core.PathDir(path)
}

func EvalSymlinks(path string) (string, error) {
	r := core.PathEvalSymlinks(path)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return "", err
		}
		return "", core.NewError(r.Error())
	}
	return r.Value.(string), nil
}

func Ext(path string) string {
	return core.PathExt(path)
}

func FromSlash(path string) string {
	return core.Replace(path, "/", string(core.PathSeparator))
}

func IsAbs(path string) bool {
	return core.PathIsAbs(path)
}

func Join(elem ...string) string {
	return core.PathJoin(elem...)
}

func Rel(basepath, targpath string) (string, error) {
	r := core.PathRel(basepath, targpath)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return "", err
		}
		return "", core.NewError(r.Error())
	}
	return r.Value.(string), nil
}

func ToSlash(path string) string {
	return core.Replace(path, string(core.PathSeparator), "/")
}

func VolumeName(path string) string {
	if len(path) >= 2 && path[1] == ':' {
		return path[:2]
	}
	return ""
}

func Walk(root string, fn WalkFunc) error {
	return core.PathWalk(root, fn)
}
