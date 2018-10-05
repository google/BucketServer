// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

var privateKey []byte
var googleAccessID string
var bucket string

func trimLeft(value string) string {
	for i := range value {
		if i > 0 {
			return value[i:]
		}
	}
	return value[:0]
}
func redirectToSignedURL(w http.ResponseWriter, r *http.Request) {
	expires := time.Now().Add(time.Second * 60)
	if gsURL, err := storage.SignedURL(bucket, trimLeft(r.URL.Path), &storage.SignedURLOptions{
		GoogleAccessID: googleAccessID,
		PrivateKey:     privateKey,
		Method:         "GET",
		Expires:        expires,
	}); err != nil {
		panic(err)
	} else {
		http.Redirect(w, r, gsURL, 302)
	}
}
func getArgument(offset int, errorMessage string, defaultValue string) string {
	if len(os.Args) <= offset {
		if len(defaultValue) > 0 {
			return defaultValue
		}
		panic(errors.New(errorMessage))
	}
	return os.Args[offset]
}
func main() {
	var (
		bucketName      = getArgument(1, "bucket name not specified", "")
		filenameKeyJSON = getArgument(2, "key JSON filename not specified", "key.json")
		portNumber      = getArgument(3, "port number not specified", "8080")
	)
	bucket = bucketName
	keyJSON, err := ioutil.ReadFile(filenameKeyJSON)
	if err != nil {
		panic(err)
	}
	if config, err := google.JWTConfigFromJSON(
		keyJSON,
		storage.ScopeReadOnly,
	); err == nil {
		googleAccessID = config.Email
		privateKey = config.PrivateKey
	} else {
		panic(err)
	}
	http.HandleFunc("/", redirectToSignedURL)
	if err := http.ListenAndServe(":"+portNumber, nil); err != nil {
		panic(err)
	}
}
