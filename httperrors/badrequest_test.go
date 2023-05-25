package httperrors

import (
	"testing"
)

func TestBadRequestAsError(t *testing.T) {
	err := BadRequest("bad_request", Field("field1", "invalid_int", "field1 should be an integer"))
	var _ error = err
	var _ HTTPCoder = err
}
