package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicTopology(t *testing.T) {
	topology := GenerateTopology(1, 1, 1, 1)

	assert.NotNil(t, topology)
	assert.Len(t, topology, 1)
	assert.Equal(t, topology["n1"][0].Name, "a1")
	assert.Equal(t, topology["n1"][0].Version, "v1")
	assert.Empty(t, topology["n1"][0].Endpoints[0].Connections)

}
