package resolver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/kerti/cloudflare-ddns/config"
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
			err:      fmt.Errorf("cannot parse IP: []"),
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
			err:      fmt.Errorf("cannot parse IP: [1.2.3.256]"),
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

	getExternalIPTestCases = []struct {
		name         string
		statusCode   int
		body         string
		resolverType string
		jsonPath     string
		result       net.IP
		errIsNil     bool
		errMsg       string
	}{
		// response not OK
		{
			name:         "responseNotOK",
			statusCode:   500,
			body:         "",
			resolverType: "text",
			jsonPath:     "",
			result:       nil,
			errIsNil:     false,
			errMsg:       "provider responded with HTTP/500",
		},
		// text resolver
		{
			name:         "emptyTextResponse",
			statusCode:   200,
			body:         "",
			resolverType: "text",
			jsonPath:     "",
			result:       nil,
			errIsNil:     false,
			errMsg:       "cannot parse IP: []",
		},
		{
			name:         "validTextResponse",
			statusCode:   200,
			body:         "1.2.3.4",
			resolverType: "text",
			jsonPath:     "",
			result:       net.ParseIP("1.2.3.4"),
			errIsNil:     true,
			errMsg:       "",
		},
		{
			name:         "invalidTextResponse",
			statusCode:   200,
			body:         "invalidResponseString",
			resolverType: "text",
			jsonPath:     "",
			result:       nil,
			errIsNil:     false,
			errMsg:       "cannot parse IP: [invalidResponseString]",
		},
		// json resolver
		{
			name:         "emptyJSONResponse",
			statusCode:   200,
			body:         "",
			resolverType: "json",
			jsonPath:     "ip",
			result:       nil,
			errIsNil:     false,
			errMsg:       "unexpected end of JSON input",
		},
		{
			name:         "validJSONResponse",
			statusCode:   200,
			body:         `{"ip": "1.2.3.4"}`,
			resolverType: "json",
			jsonPath:     "ip",
			result:       net.ParseIP("1.2.3.4"),
			errIsNil:     true,
			errMsg:       "",
		},
		{
			name:         "invalidJSONResponse",
			statusCode:   200,
			body:         `{"ip": "invalidResponseString"}`,
			resolverType: "json",
			jsonPath:     "ip",
			result:       nil,
			errIsNil:     false,
			errMsg:       "cannot parse IP: [invalidResponseString]",
		},
		// unsupported resolver type
		{
			name:         "unsupportedResolverType",
			statusCode:   200,
			body:         `{"ip": "invalidResponseString"}`,
			resolverType: "random",
			jsonPath:     "ip",
			result:       nil,
			errIsNil:     false,
			errMsg:       "unsupported resolver type: [random]",
		},
	}
)

// RoundTripFunc is the signature for fake transport func
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is the fake round tripper for our fake HTTP client
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func mockHTTPResponse(input string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(input)),
	}
}

func TestResolver(t *testing.T) {

	logger.InitLogger(&initloglevel)
	configFile := "../config.yaml"
	config.Load(&configFile)

	t.Run("Get", func(t *testing.T) {
		res, err := Get()
		assert.NotNil(t, res)
		assert.NotEmpty(t, res)
		assert.Nil(t, err)
	})

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

	t.Run("GetExternalIP", func(t *testing.T) {

		for _, tc := range getExternalIPTestCases {
			t.Run(tc.name, func(t *testing.T) {
				client := NewTestClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       ioutil.NopCloser(bytes.NewBufferString(tc.body)),
						Header:     make(http.Header),
					}
				})

				resolver := Resolver{
					HTTPClient: client,
					Type:       tc.resolverType,
					JSONPath:   tc.jsonPath,
				}

				res, err := resolver.GetExternalIP()
				if tc.errIsNil {
					assert.Nil(t, err)
					assert.NotNil(t, res)
					assert.Equal(t, tc.result, res)
				} else {
					assert.NotNil(t, err)
					assert.Equal(t, tc.errMsg, err.Error())
					assert.Nil(t, res)
				}
			})
		}
	})
}
