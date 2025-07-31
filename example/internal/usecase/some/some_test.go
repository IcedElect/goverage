package some

import (
	"testing"
)

func TestSumMethod(t *testing.T) {
	u := NewUsecase()

	result := u.SumMethod(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
