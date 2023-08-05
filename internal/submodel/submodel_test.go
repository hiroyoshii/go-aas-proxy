package submodel

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	os.Setenv("SUBMODEL_CONFIG_PATH", "../../test/configs/submodel_config.yaml")
	s, err := NewSubmodel()

	assert.Nil(t, err)
	assert.NotNil(t, s)
	b, err := s.Get("test", "test")
	assert.Nil(t, err)
	assert.Nil(t, string(b))
}
