package resolver

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/stretchr/testify/assert"
)

var (
	initloglevel uint8 = uint8(0)

	readIPTextTestCases = []struct {
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

	readIPJSONTestCases = []struct {
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
			errMsg:   "cannot parse IP: [somerandomtext]",
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

func mockHTTPResponse(input string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(input)),
	}
}

func TestResolver(t *testing.T) {

	logger.InitLogger(&initloglevel)

	t.Run("init", func(t *testing.T) {
		resolver := Resolver{}
		resolver.Init()
		assert.NotNil(t, resolver.HTTPClient)
	})

	t.Run("readIPText", func(t *testing.T) {

		resolver := Resolver{
			Type: "text",
			URL:  "",
		}

		for _, tc := range readIPTextTestCases {
			t.Run(tc.name, func(t *testing.T) {
				mockedResponse := mockHTTPResponse(tc.response)
				result, err := resolver.readIPText(*mockedResponse)
				assert.Equal(t, tc.result, result)
				assert.Equal(t, tc.err, err)
			})
		}

	})

	t.Run("readIPJSON", func(t *testing.T) {

		resolver := Resolver{
			Type: "json",
			URL:  "",
		}

		for _, tc := range readIPJSONTestCases {
			t.Run(tc.name, func(t *testing.T) {
				resolver.JSONPath = tc.jsonPath
				mockedResponse := mockHTTPResponse(tc.response)
				result, err := resolver.readIPJSON(*mockedResponse)
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
