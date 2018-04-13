package rpc

import "testing"

func TestGetBlockByHeight_EmptyUrl(t *testing.T) {
	_, err := GetBlockByHeight(0)
	if err == nil || err.Error() != "Unknown rpc url." {
		t.Error("Expect an error")
	}
}

func TestGetCurrentHeight_EmptyUrl(t *testing.T) {
	_, err := GetCurrentHeight()
	if err == nil || err.Error() != "Unknown rpc url." {
		t.Error("Expect an error")
	}
}
