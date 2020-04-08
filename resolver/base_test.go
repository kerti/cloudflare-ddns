package resolver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/stretchr/testify/assert"
)

var (
	// initial loglevel for mocking logger
	initloglevel uint8 = uint8(0)
)

func getSimpleResolver(url string) *Resolver {
	var resolver Resolver
	resolver = &Simple{URL: url}
	return &resolver
}

func getJSONResolver(url string) *Resolver {
	var resolver Resolver
	resolver = &JSON{URL: url}
	return &resolver
}

func mockHTTPResponse(input string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(input)),
	}
}

func TestBase(t *testing.T) {

	logger.InitLogger(&initloglevel)

	t.Run("Get", func(t *testing.T) {
		testCases := []struct {
			name   string
			key    string
			result *Resolver
			err    error
		}{
			{
				name:   "emptyKey",
				key:    "",
				result: nil,
				err:    fmt.Errorf("unsupported resolver: "),
			},
			{
				name:   "randomKey",
				key:    "someRandomText",
				result: nil,
				err:    fmt.Errorf("unsupported resolver: someRandomText"),
			},
			{
				name:   "validJSONResolverKey",
				key:    ResolverBigDataCloud,
				result: getSimpleResolver(ResolverURLs[ResolverBigDataCloud]),
				err:    nil,
			},
			{
				name:   "validSimpleResolverKey",
				key:    ResolverICanHazIP,
				result: getJSONResolver(ResolverURLs[ResolverICanHazIP]),
				err:    nil,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := Get(tc.key)
				assert.Equal(t, tc.result == nil, result == nil)
				assert.Equal(t, tc.err, err)
			})
		}
	})

	t.Run("GetAll", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			res, err := GetAll()
			assert.Nil(t, err)
			assert.Equal(t, len(ResolverURLs), len(res))
		})

		t.Run("error", func(t *testing.T) {
			key := "randomString"
			ResolverURLs[key] = key

			res, err := GetAll()

			assert.Equal(t, fmt.Errorf("unsupported resolver: %v", key), err)
			assert.Nil(t, res)

			delete(ResolverURLs, key)
		})

	})
}
