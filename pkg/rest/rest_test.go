package rest

import (
	"reflect"
	"testing"
)

const (
	tstMessage = "test"
)

type TestStruct struct {
}

func TestCreateResponse(t *testing.T) {
	var payloadType reflect.Type
	var testPayload TestStruct
	target := &REST[any]{
		data: testPayload,
		formatter: func(payload any) string {
			payloadType = reflect.TypeOf(payload)
			return tstMessage
		},
	}

	res := target.createReponse()

	if res.Data.Content != tstMessage || payloadType != reflect.TypeOf(testPayload) {
		t.Error("Reponse message does not match!")
	}
}
