package json

import core "dappco.re/go"

func Marshal(v any) ([]byte, error) {
	r := core.JSONMarshal(v)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, core.NewError(r.Error())
	}
	return r.Value.([]byte), nil
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	r := core.JSONMarshalIndent(v, prefix, indent)
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, core.NewError(r.Error())
	}
	return r.Value.([]byte), nil
}

func Unmarshal(data []byte, target any) error {
	r := core.JSONUnmarshal(data, target)
	if r.OK {
		return nil
	}
	if err, ok := r.Value.(error); ok {
		return err
	}
	return core.NewError(r.Error())
}
