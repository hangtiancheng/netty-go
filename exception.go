package netty

import "fmt"

// Exception defines an exception
type Exception = error

func AsException(ex any) Exception {
	if nil != ex {
		if e, ok := ex.(Exception); ok {
			return e
		}
		return fmt.Errorf("%v", ex)
	}
	return nil
}
