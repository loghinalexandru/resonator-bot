package rest

import (
	"testing"
)

const (
	tstMessage = "test"
)

type TestStruct struct {
}

func TestCreateResponse(t *testing.T) {
	var testPayload TestStruct
	target := &REST[TestStruct]{
		data: testPayload,
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	res := target.createReponse()

	if res.Data.Content != tstMessage {
		t.Error("Reponse message does not match!")
	}
}
