package generator

import (
	"testing"

	"github.com/leandroberetta/mimik/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestBasicTopology(t *testing.T) {
	generator := api.Generator{
		Services:          1,
		Namespaces:        1,
		Connections:       1,
		RandomConnections: 1,
	}
	topology := GenerateTopology(generator)

	assert.NotNil(t, topology)
	assert.Len(t, topology, 1)
	assert.Equal(t, topology["n1"][0].Name, "a1")
	assert.Equal(t, topology["n1"][0].Version, "v1")
	assert.Empty(t, topology["n1"][0].Endpoints[0].Connections)

}
