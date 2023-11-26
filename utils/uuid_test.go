package utils

import "testing"

func TestCreateUUID(t *testing.T) {

	for i := 0; i < 10; i++ {
		t.Log(CreateGUID())
	}
}
