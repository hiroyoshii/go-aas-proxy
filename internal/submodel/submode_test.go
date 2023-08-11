package submodel

import (
	"encoding/json"
	basyxAas "hiroyoshii/go-aas-proxy/gen/go"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubmodel(t *testing.T) {
	os.Setenv("SUBMODEL_CONFIG_PATH", "submodel_config.yaml")

	s, err := NewSubmodel()
	assert.Nil(t, err)
	b, err := s.Get("454576463545648365874", "https://www.hsu-hh.de/aut/aas/identification", "Identification")
	assert.Nil(t, err)

	var res *basyxAas.Submodel
	err = json.Unmarshal(b, &res)
	assert.Nil(t, err)

}
