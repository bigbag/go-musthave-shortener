package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/bigbag/go-musthave-shortener/internal/config"
)

var (
	testServer *Server
)

type TestCase struct {
	description    string
	requestRoute   string
	requestMethod  string
	requestBody    string
	expectedError  bool
	expectedCode   int
	expectedBody   string
	requestHeaders map[string][]string
}

func getNewTestServer() *Server {
	if testServer == nil {
		cfg, _ := config.New()
		testServer = New(logrus.New(), cfg)
	}

	return testServer
}

func makeTestRequest(server *Server, test TestCase) (*http.Response, error) {
	req, _ := http.NewRequest(
		test.requestMethod,
		test.requestRoute,
		strings.NewReader(test.requestBody),
	)
	if test.requestHeaders != nil {
		req.Header = test.requestHeaders
	}

	return server.f.Test(req, -1)

}
func checkResponse(t *testing.T, test TestCase, res *http.Response, err error) {
	assert.Equalf(t, test.expectedError, err != nil, test.description)

	if test.expectedError {
		return
	}

	assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

	body, err := ioutil.ReadAll(res.Body)

	assert.Nilf(t, err, test.description)

	if test.expectedBody == "" {
		return
	}
	assert.Equalf(t, test.expectedBody, string(body), test.description)
}

func TestGetStatusHandler(t *testing.T) {
	tests := []TestCase{
		{
			description:   "success",
			requestRoute:  "/ping",
			requestMethod: http.MethodGet,
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"result":"OK"}`,
		},
	}

	server := getNewTestServer()
	for _, test := range tests {
		res, err := makeTestRequest(server, test)
		checkResponse(t, test, res, err)
	}
}

func TestCreateURLHandler(t *testing.T) {
	tests := []TestCase{
		{
			description:   "method not allowed",
			requestRoute:  "/",
			requestMethod: http.MethodGet,
			expectedError: false,
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"code":500,"message":"Method Not Allowed"}`,
		},
		{
			description:   "empty payload",
			requestRoute:  "/",
			requestMethod: http.MethodPost,
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"code":400,"message":"Please specify a valid full url"}`,
		},
		{
			description:   "success",
			requestRoute:  "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://github.com/plain",
			expectedError: false,
			expectedCode:  http.StatusCreated,
			expectedBody:  "",
		},
		{
			description:   "conflict",
			requestRoute:  "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://github.com/plain",
			expectedError: false,
			expectedCode:  http.StatusConflict,
			expectedBody:  "",
		},
	}

	server := getNewTestServer()
	for _, test := range tests {
		res, err := makeTestRequest(server, test)
		checkResponse(t, test, res, err)
	}
}

func TestFullFlow(t *testing.T) {
	originFullURL := "https://github.com/full_plain"
	server := getNewTestServer()

	test := TestCase{
		description:   "create url",
		requestRoute:  "/",
		requestMethod: http.MethodPost,
		requestBody:   originFullURL,
		expectedError: false,
		expectedCode:  http.StatusCreated,
		expectedBody:  "",
	}
	res, err := makeTestRequest(server, test)
	body, _ := ioutil.ReadAll(res.Body)
	shortURL, _ := url.Parse(string(body))

	checkResponse(t, test, res, err)

	test = TestCase{
		description:   "get full url",
		requestRoute:  shortURL.Path,
		requestMethod: http.MethodGet,
		expectedError: false,
		expectedCode:  http.StatusTemporaryRedirect,
		expectedBody:  "",
	}
	res, err = makeTestRequest(server, test)
	fullURL := res.Header.Values("Location")[0]

	checkResponse(t, test, res, err)

	assert.Equalf(t, originFullURL, fullURL, test.description)

}

func TestCreateURLJsonHandler(t *testing.T) {
	tests := []TestCase{
		{
			description:   "method not allowed",
			requestRoute:  "/api/shorten",
			requestMethod: http.MethodGet,
			expectedError: false,
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"code":500,"message":"Method Not Allowed"}`,
		},
		{
			description:   "empty payload",
			requestRoute:  "/api/shorten",
			requestMethod: http.MethodPost,
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"code":400,"message":"Please specify a valid full url"}`,
		},
		{
			description:   "not valid content type",
			requestRoute:  "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url":"https://github.com/json"}`,
			requestHeaders: http.Header{
				"Content-Type": []string{"text/html"},
			},
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"code":400,"message":"Please specify a valid full url"}`,
		},
		{
			description:   "success",
			requestRoute:  "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url":"https://github.com/json"}`,
			requestHeaders: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedError: false,
			expectedCode:  http.StatusCreated,
			expectedBody:  "",
		},
		{
			description:   "conflict",
			requestRoute:  "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url":"https://github.com/json"}`,
			requestHeaders: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedError: false,
			expectedCode:  http.StatusConflict,
			expectedBody:  "",
		},
	}

	server := getNewTestServer()
	for _, test := range tests {
		res, err := makeTestRequest(server, test)
		checkResponse(t, test, res, err)
	}
}

func TestFullFlowJson(t *testing.T) {
	type ShorterResult struct {
		Result string `json:"result"`
	}

	originFullURL := "https://github.com/full_json"
	server := getNewTestServer()

	test := TestCase{
		description:   "create url",
		requestRoute:  "/api/shorten",
		requestMethod: http.MethodPost,
		requestBody:   fmt.Sprintf(`{"url":"%s"}`, originFullURL),
		requestHeaders: http.Header{
			"Content-Type": []string{"application/json"},
		},
		expectedError: false,
		expectedCode:  http.StatusCreated,
		expectedBody:  "",
	}
	res, err := makeTestRequest(server, test)

	var shorterResult ShorterResult
	json.NewDecoder(res.Body).Decode(&shorterResult)
	shortURL, _ := url.Parse(shorterResult.Result)

	checkResponse(t, test, res, err)

	test = TestCase{
		description:   "get full url",
		requestRoute:  shortURL.Path,
		requestMethod: http.MethodGet,
		expectedError: false,
		expectedCode:  http.StatusTemporaryRedirect,
		expectedBody:  "",
	}
	res, err = makeTestRequest(server, test)
	fullURL := res.Header.Values("Location")[0]

	checkResponse(t, test, res, err)

	assert.Equalf(t, originFullURL, fullURL, test.description)

}

func TestCreateBatchOfShortURLHandler(t *testing.T) {
	type ShorterItem struct {
		ShortURL      string `json:"short_url"`
		CorrelationID string `json:"correlation_id"`
	}

	type ShorterResult []*ShorterItem

	originFullURL := "https://github.com/batch_json"
	server := getNewTestServer()

	requestJSON := `[
						{
							"original_url":"%s",
							"correlation_id": "323"
						}
					]`

	test := TestCase{
		description:   "create url",
		requestRoute:  "/api/shorten/batch",
		requestMethod: http.MethodPost,
		requestBody:   fmt.Sprintf(requestJSON, originFullURL),
		requestHeaders: http.Header{
			"Content-Type": []string{"application/json"},
		},
		expectedError: false,
		expectedCode:  http.StatusCreated,
		expectedBody:  "",
	}
	res, err := makeTestRequest(server, test)

	var shorterResult ShorterResult
	json.NewDecoder(res.Body).Decode(&shorterResult)
	shortURL, _ := url.Parse(shorterResult[0].ShortURL)

	checkResponse(t, test, res, err)

	test = TestCase{
		description:   "get full url",
		requestRoute:  shortURL.Path,
		requestMethod: http.MethodGet,
		expectedError: false,
		expectedCode:  http.StatusTemporaryRedirect,
		expectedBody:  "",
	}
	res, err = makeTestRequest(server, test)
	fullURL := res.Header.Values("Location")[0]

	checkResponse(t, test, res, err)

	assert.Equalf(t, originFullURL, fullURL, test.description)

}
