package tests

import (
	. "github.com/advancedlogic/goms/pkg/models"
	. "gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestMicroservice(t *testing.T) {
	_, err := NewMicroservice()
	Equal(t, err, nil)
}
