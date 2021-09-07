package eyeson

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

const testAPIKey = "secret-api-key"

// setup sets up a test HTTP server along with a eyeson.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
// Ref: https://github.com/google/go-github/blob/master/github/github_test.go
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()
	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	client = NewClient(testAPIKey)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

// testMethod provides a helper function to test a request method type is as
// expected.
// Ref: https://github.com/google/go-github/blob/master/github/github_test.go
func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

// testFormValues provides a helper function to test if request parameters are
// provided as expected.
// Ref: https://github.com/google/go-github/blob/master/github/github_test.go
func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("apiKey")

	if got, want := c.BaseURL.String(), endpoint; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.apiKey, "apiKey"; got != want {
		t.Errorf("NewClient apiKey is %v, want %v", got, want)
	}
}

func TestNewRequest_authorization(t *testing.T) {
	apiKey := "secret-key"
	c := NewClient(apiKey)
	req, err := c.NewRequest("GET", ".", nil)

	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if got, want := req.Header.Get("Authorization"), apiKey; got != want {
		t.Fatalf("NewRequest apiKey is %v, want %v", got, want)
	}
}

func TestNewRequest_userAuthorization(t *testing.T) {
	c := NewClient("api-key").UserClient()
	req, err := c.NewRequest("GET", ".", nil)

	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization header was set (%v) but should not", got)
	}
}

func TestNewRequest_userAgent(t *testing.T) {
	c := NewClient("")
	req, err := c.NewRequest("GET", ".", nil)

	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if got, want := req.Header.Get("User-Agent"), "eyeson-go"; got != want {
		t.Fatalf("NewRequest userAgent is %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type testRoom struct {
		Title string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"title":"foobar"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(testRoom)
	client.Do(req, body)

	want := &testRoom{"foobar"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Do body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})
	return

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp != nil {
		t.Errorf("Expected empty response, got %v", resp.Body)
	}
}
