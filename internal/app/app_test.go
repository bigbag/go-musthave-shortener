package app

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/bigbag/go-musthave-shortener/internal/config"
)

type TestCase struct {
	description   string
	requestRoute  string
	requestMethod string
	requestBody   string
	expectedError bool
	expectedCode  int
	expectedBody  string
}

func getNewTestServer() *Server {
	cfg, _ := config.New()
	return New(logrus.New(), cfg)
}

func makeTestRequest(server *Server, test TestCase) (*http.Response, error) {
	req, _ := http.NewRequest(
		test.requestMethod,
		test.requestRoute,
		strings.NewReader(test.requestBody),
	)

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

func TestCreateURLHandler(t *testing.T) {
	tests := []TestCase{
		{
			description:   "method not allowed",
			requestRoute:  "/",
			requestMethod: "GET",
			requestBody:   "",
			expectedError: false,
			expectedCode:  500,
			expectedBody:  `{"code":500,"message":"Method Not Allowed"}`,
		},
		{
			description:   "empty payload",
			requestRoute:  "/",
			requestMethod: "POST",
			requestBody:   "",
			expectedError: false,
			expectedCode:  400,
			expectedBody:  `{"code":400,"message":"Please specify a valid full url"}`,
		},
		{
			description:   "success",
			requestRoute:  "/",
			requestMethod: "POST",
			requestBody:   "https://github.com",
			expectedError: false,
			expectedCode:  201,
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
	originFullURL := "https://github.com"
	server := getNewTestServer()

	test := TestCase{
		description:   "create url",
		requestRoute:  "/",
		requestMethod: "POST",
		requestBody:   originFullURL,
		expectedError: false,
		expectedCode:  201,
		expectedBody:  "",
	}
	res, err := makeTestRequest(server, test)
	body, _ := ioutil.ReadAll(res.Body)
	shortURL, _ := url.Parse(string(body))

	checkResponse(t, test, res, err)

	test = TestCase{
		description:   "get full url",
		requestRoute:  shortURL.Path,
		requestMethod: "GET",
		requestBody:   "",
		expectedError: false,
		expectedCode:  307,
		expectedBody:  "",
	}
	res, err = makeTestRequest(server, test)
	fullURL := res.Header.Values("Location")[0]

	checkResponse(t, test, res, err)

	assert.Equalf(t, originFullURL, fullURL, test.description)

}
