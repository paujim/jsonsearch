package jsonsearch

import (
	"testing"
)

func TestIsEqualTo(t *testing.T) {

	obj := NewJsonObject(nil, test2, "principal", "id")
	if obj == nil {
		t.Errorf("expected: a valid json object")
	}
	otherObj := NewJsonObject(nil, test1, "other", "id")
	if otherObj == nil {
		t.Errorf("expected: a valid json object")
	}

	actualResponse := obj.Where("active").IsEqualTo("true")

	if len(actualResponse) != 2 {
		t.Errorf("expected: [active == true] should have 2 values in the response")
	}

	actualResponse = obj.Where("created_at").Contains("2016")

	if len(actualResponse) != 3 {
		t.Errorf("expected: [created_at == 2016] should have 3 values in the response")
	}

	actualResponse = obj.With(map[string]*JsonObject{"other_id": otherObj}).Where("created_at").Contains("2016")

	if len(actualResponse) != 3 {
		t.Errorf("expected: [created_at == 2016] should have 3 values in the response")
	}
}
