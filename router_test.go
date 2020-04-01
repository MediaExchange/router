/*
   Copyright 2019 MediaExchange.io

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package router

import (
	"fmt"
	"github.com/MediaExchange/assert"
	"github.com/MediaExchange/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRouter() *Router {
	log.SetWriter(ioutil.Discard)
	return NewRouter()
}

func TestNewRouter(t *testing.T) {
	r := NewRouter()

	assert.
		With(t).
		That(r).
		IsNotNil()
}

func TestRouter_AddRoute(t *testing.T) {
	r := newRouter()
	r.AddRoute("GET", "/path/{foo}/{bar}", func(writer http.ResponseWriter, request *http.Request) {
	})

	Assert := assert.With(t)
	Assert.
		That(r.routes).
		IsNotNil()
	Assert.
		That(len(r.routes)).
		IsEqualTo(1)

	s := r.routes[0]
	Assert.
		That(s.handler).
		IsNotNil()

	actual := s.path.String()
	expected := "^/path/(?P<foo>.+?)/(?P<bar>.+?)$"

	Assert.
		That(actual).
		IsEqualTo(expected)
	Assert.
		That(s.verb).
		IsEqualTo("GET")
}

func TestRouter_ServeHTTP(t *testing.T) {
	Assert := assert.With(t)

	// Create the router
	router := newRouter()
	router.AddRoute("GET", "/path/{foo}/{bar}", func(writer http.ResponseWriter, request *http.Request) {
		params := GetParams(request.Context())

		Assert.
			That(len(params)).
			IsEqualTo(4)
		Assert.
			That(params["foo"]).
			IsEqualTo("fuz")
		Assert.
			That(params["bar"]).
			IsEqualTo("baz")
		Assert.
			That(params["aaa"]).
			IsEqualTo("bbb")
		Assert.
			That(params["ccc"]).
			IsEqualTo("ddd")

		writer.WriteHeader(200)
	})

	router.AddRoute("GET", "/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	request := server.URL + "/path/fuz/baz?aaa=bbb&ccc=ddd"
	res, err := http.Get(request)

	Assert.
		That(err).
		IsOk()
	Assert.
		That(res.StatusCode).
		IsEqualTo(200)

	request = server.URL + "/"
	res, err = http.Get(request)

	Assert.
		That(err).
		IsOk()
	Assert.
		That(res.StatusCode).
		IsEqualTo(200)
}

// Examples that describe how to implement route handlers that accept query
// parameters and path variables.
func Example() {
	// Sample search handler.
	search := func(writer http.ResponseWriter, request *http.Request) {
		params := GetParams(request.Context())
		fmt.Printf("Search for: \"%s\"\n", params["s"])
		writer.WriteHeader(200)
	}

	// Sample get resource handler.
	getBookByIsbn := func(writer http.ResponseWriter, request *http.Request) {
		params := GetParams(request.Context())
		fmt.Printf("Get book with ISBN = %s\n", params["isbn"])
		writer.WriteHeader(200)
	}

	// build the router.
	router := NewRouter().
		AddRoute("GET", "/search", search).
		AddRoute("GET", "/book/{isbn}", getBookByIsbn)

	// Start the server.
	server := httptest.NewServer(router)
	defer server.Close()

	http.Get(server.URL + "/search?s=stuff+to+search+for")
	http.Get(server.URL + "/book/978-0316371247")

	// OUTPUT: Search for: "stuff to search for"
	// Get book with ISBN = 978-0316371247
}
