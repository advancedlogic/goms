package tests

import (
	. "github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
	. "gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestMicroservice(t *testing.T) {
	env, err := modules.NewEnvironmentBuilder().WithConfig("config").WithName("micro").Build()
	Equal(t, err, nil)
	_, err = NewMicroservice(env)
	Equal(t, err, nil)
}
