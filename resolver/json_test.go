package resolver

import (
	"net"
	"testing"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/stretchr/testify/assert"
)

var (
	// test vars go here
	jsonReadIPTestCases = []struct {
		name     string
		jsonPath string
		response string
		result   net.IP
		errIsNil bool
		errMsg   string
	}{
		{
			name:     "emptyResponse",
			jsonPath: "path",
			response: "",
			result:   nil,
			errIsNil: false,
			errMsg:   "unexpected end of JSON input",
		},
		{
			name:     "plaintextResponse",
			jsonPath: "path",
			response: "1.2.3.4",
			result:   nil,
			errIsNil: false,
			errMsg:   "invalid character '.' after top-level value",
		},
		{
			name:     "validSimpleJSONResponseWithInvalidPath",
			jsonPath: "invalidPath",
			response: `{"ipAddress": "1.2.3.4"}`,
			result:   nil,
			errIsNil: false,
			errMsg:   "IP address not found at path [invalidPath]",
		},
		{
			name:     "validSimpleJSONResponseWithInvalidValueType",
			jsonPath: "ipAddress",
			response: `{"ipAddress": true}`,
			result:   nil,
			errIsNil: false,
			errMsg:   "cannot convert value at path [ipAddress] to string: true",
		},
		{
			name:     "validSimpleJSONResponseWithInvalidValue",
			jsonPath: "ipAddress",
			response: `{"ipAddress": "somerandomtext"}`,
			result:   nil,
			errIsNil: false,
			errMsg:   "cannot parse IP: somerandomtext",
		},
		{
			name:     "validSimpleJSONResponse",
			jsonPath: "ipAddress",
			response: `{"ipAddress": "1.2.3.4"}`,
			result:   net.ParseIP("1.2.3.4"),
			errIsNil: true,
			errMsg:   "",
		},
	}
)

func TestJSONResolver(t *testing.T) {

	logger.InitLogger(&initloglevel)

	t.Run("readIP", func(t *testing.T) {

		resolver := JSON{
			URL: "",
		}

		t.Run("nilResponse", func(t *testing.T) {
			result, err := resolver.readIP(nil)
			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.Equal(t, "response is nil", err.Error())
		})

		for _, tc := range jsonReadIPTestCases {
			t.Run(tc.name, func(t *testing.T) {
				resolver.JSONPath = tc.jsonPath
				result, err := resolver.readIP(mockHTTPResponse(tc.response))
				assert.Equal(t, tc.result, result)
				if tc.errIsNil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			})
		}
	})
}
