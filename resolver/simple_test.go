package resolver

import (
	"fmt"
	"net"
	"testing"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/stretchr/testify/assert"
)

var (
	// test vars go here
	simpleReadIPTestCases = []struct {
		name     string
		response string
		result   net.IP
		err      error
	}{
		{
			name:     "emptyResponse",
			response: "",
			result:   nil,
			err:      fmt.Errorf("cannot parse IP: "),
		},
		{
			name:     "validIpAddress",
			response: "1.2.3.4",
			result:   net.ParseIP("1.2.3.4"),
			err:      nil,
		},
		{
			name:     "invalidIpAddress",
			response: "1.2.3.256",
			result:   nil,
			err:      fmt.Errorf("cannot parse IP: 1.2.3.256"),
		},
		{
			name:     "validIpAddressWithNewline",
			response: "1.2.3.4\n",
			result:   net.ParseIP("1.2.3.4"),
			err:      nil,
		},
		{
			name:     "validIpAddressWithQuotes",
			response: `"1.2.3.4"`,
			result:   net.ParseIP("1.2.3.4"),
			err:      nil,
		},
	}
)

func TestSimpleResolver(t *testing.T) {

	logger.InitLogger(&initloglevel)

	t.Run("readIP", func(t *testing.T) {

		resolver := Simple{
			URL: "",
		}

		t.Run("nilResponse", func(t *testing.T) {
			result, err := resolver.readIP(nil)
			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.Equal(t, "response is nil", err.Error())
		})

		for _, tc := range simpleReadIPTestCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := resolver.readIP(mockHTTPResponse(tc.response))
				assert.Equal(t, tc.result, result)
				assert.Equal(t, tc.err, err)
			})
		}
	})
}
