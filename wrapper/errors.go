package wrapper

import (
	"fmt"
)

var ErrUnknownRequestDataType = func(data any) error {
	return fmt.Errorf("unknown request data type. type: %#v", data)
}
