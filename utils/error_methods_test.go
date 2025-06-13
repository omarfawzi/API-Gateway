package utils

import "testing"

func TestHTTPResponseErrorMethods(t *testing.T) {
	errObj := HTTPResponseError{Code: 404, Msg: "not found", HTTPEncoding: "app"}
	if errObj.Error() != "not found" {
		t.Errorf("Error() = %s", errObj.Error())
	}
	if errObj.StatusCode() != 404 {
		t.Errorf("StatusCode() = %d", errObj.StatusCode())
	}
	if errObj.Encoding() != "app" {
		t.Errorf("Encoding() = %s", errObj.Encoding())
	}
}
