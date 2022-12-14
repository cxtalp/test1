/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2017 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package signer

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

// Tests url encoding.
func TestEncodeURL2Path(t *testing.T) {
	type urlStrings struct {
		virtualHost    bool
		bucketName     string
		objName        string
		encodedObjName string
	}

	want := []urlStrings{
		{
			virtualHost:    true,
			bucketName:     "bucketName",
			objName:        "本語",
			encodedObjName: "%E6%9C%AC%E8%AA%9E",
		},
		{
			virtualHost:    true,
			bucketName:     "bucketName",
			objName:        "本語.1",
			encodedObjName: "%E6%9C%AC%E8%AA%9E.1",
		},
		{
			virtualHost:    true,
			objName:        ">123>3123123",
			bucketName:     "bucketName",
			encodedObjName: "%3E123%3E3123123",
		},
		{
			virtualHost:    true,
			bucketName:     "bucketName",
			objName:        "test 1 2.txt",
			encodedObjName: "test%201%202.txt",
		},
		{
			virtualHost:    false,
			bucketName:     "test.bucketName",
			objName:        "test++ 1.txt",
			encodedObjName: "test%2B%2B%201.txt",
		},
	}

	for i, o := range want {
		var hostURL string
		if o.virtualHost {
			hostURL = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", o.bucketName, o.objName)
		} else {
			hostURL = fmt.Sprintf("https://s3.amazonaws.com/%s/%s", o.bucketName, o.objName)
		}
		u, err := url.Parse(hostURL)
		if err != nil {
			t.Fatalf("Test %d, Error: %v", i+1, err)
		}
		expectedPath := "/" + o.bucketName + "/" + o.encodedObjName
		if foundPath := encodeURL2Path(&http.Request{URL: u}, o.virtualHost); foundPath != expectedPath {
			t.Fatalf("Test %d, Error: expected = `%v`, found = `%v`", i+1, expectedPath, foundPath)
		}
	}
}

// TestSignV4TrimAll - tests the logic of TrimAll() function
func TestSignV4TrimAll(t *testing.T) {
	testCases := []struct {
		// Input.
		inputStr string
		// Expected result.
		result string
	}{
		{"本語", "本語"},
		{" abc ", "abc"},
		{" a b ", "a b"},
		{"a b ", "a b"},
		{"a  b", "a b"},
		{"a   b", "a b"},
		{"   a   b  c   ", "a b c"},
		{"a \t b  c   ", "a b c"},
		{"\"a \t b  c   ", "\"a b c"},
		{" \t\n\u000b\r\fa \t\n\u000b\r\f b \t\n\u000b\r\f c \t\n\u000b\r\f", "a b c"},
	}

	// Tests generated values from url encoded name.
	for i, testCase := range testCases {
		result := signV4TrimAll(testCase.inputStr)
		if testCase.result != result {
			t.Errorf("Test %d: Expected signV4TrimAll result to be \"%s\", but found it to be \"%s\" instead", i+1, testCase.result, result)
		}
	}
}
