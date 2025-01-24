package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	var hello Hello
	hello.Type = RC4
	var buf, _ = hello.ToBytes()
	assert.Equal(t, 1, len(buf), "to bytes ok")
	assert.Equal(t, nil, hello.FromBytes(buf), "from byte ok")
	assert.Equal(t, hello.Type, RC4, "value ok")
}
