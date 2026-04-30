package application

import core "dappco.re/go"

func jsonMarshal(value any) ([]byte, resultFailure) {
	result := core.JSONMarshal(value)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, err
		}
		return nil, core.NewError(result.Error())
	}
	return result.Value.([]byte), nil
}

func jsonUnmarshal(data []byte, target any) resultFailure {
	result := core.JSONUnmarshal(data, target)
	if result.OK {
		return nil
	}
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.NewError(result.Error())
}
