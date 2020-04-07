package resolver

import (
	"testing"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/stretchr/testify/assert"
)

var (
	// initial loglevel for mocking logger
	initloglevel uint8 = uint8(0)
)

func TestBase(t *testing.T) {

	logger.InitLogger(&initloglevel)

	t.Run("GetAll", func(t *testing.T) {
		res, err := GetAll()
		assert.Equal(t, nil, err)
		assert.Equal(t, len(ResolverURLs), len(res))
	})
}
