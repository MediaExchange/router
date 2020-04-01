# router

[![Build Status](https://travis-ci.org/mediaexchange/router.svg)](https://travis-ci.org/mediaexchange/router)
[![GoDoc](https://godoc.org/github.com/mediaexchange/router/github?status.svg)](https://godoc.org/github.com/mediaexchange/router)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go version](https://img.shields.io/badge/go-~%3E1.13-green.svg)](https://golang.org/doc/devel/release.html#go1.13)
[![Go version](https://img.shields.io/badge/go-~%3E1.14-green.svg)](https://golang.org/doc/devel/release.html#go1.14)

This is a simple HTTP router library that does not contain any middleware,
but that doesn't prevent the use of other middleware libraries. Handler
functions are identical to the normal `ServeHTTP(http.ResponseWriter, *http.Request)`
but the `request.Context()` object contains a `map` of the path variables and
query parameters in a field named `params`. 

## Features

* Uses the same machinery as the standard HTTP server.
* Should integrate well with all available middleware libraries.
* Decodes path variables and query parameters.
* `AddRoute` function can be chained.
* `GetParams` helper function to extract the path parameters and query
    parameters from the request context. 

## Usage

```go
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
```

## Contributing

 1.  Fork it
 2.  Create a feature branch (`git checkout -b new-feature`)
 3.  Commit changes (`git commit -am "Added new feature xyz"`)
 4.  Push the branch (`git push origin new-feature`)
 5.  Create a new pull request.

## Maintainers

* [Media Exchange](http://github.com/MediaExchange)

## License

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
