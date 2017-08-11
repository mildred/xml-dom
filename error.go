package xmldom

import (
	"fmt"
)

type DOMError struct {
	code ErrorCode
}

func (e *DOMError) Error() string {
	switch e.code {
	case IndexSizeError:
		return fmt.Sprintf("Error, index size")
	case DOMStringSizeError:
		return fmt.Sprintf("Error, string size")
	case HierarchyRequestError:
		return fmt.Sprintf("Error, hierarchy request")
	case WrongDocumentError:
		return fmt.Sprintf("Error, wrong document")
	case InvalidCharacterError:
		return fmt.Sprintf("Error, invalid character")
	case NoDataAllowedError:
		return fmt.Sprintf("Error, no data allowed")
	case NoModificationAllowedError:
		return fmt.Sprintf("Error, no modification allowed")
	case NotFoundError:
		return fmt.Sprintf("Error, not found")
	case NotSupportedError:
		return fmt.Sprintf("Error, not supported")
	case InuseAttributeError:
		return fmt.Sprintf("Error, node already in use")
	default:
		return fmt.Sprintf("Error code %d", e.code)
	}
}

func (e *DOMError) Code() ErrorCode {
	return e.code
}

func err(c ErrorCode) *DOMError {
	return &DOMError{c}
}
