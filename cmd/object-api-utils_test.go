/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
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

package cmd

import (
	"io/ioutil"
	"strings"
	"testing"
)

// Tests validate bucket name.
func TestIsValidBucketName(t *testing.T) {
	testCases := []struct {
		bucketName string
		shouldPass bool
	}{
		// cases which should pass the test.
		// passing in valid bucket names.
		{"lol", true},
		{"1-this-is-valid", true},
		{"1-this-too-is-valid-1", true},
		{"this.works.too.1", true},
		{"1234567", true},
		{"123", true},
		{"s3-eu-west-1.amazonaws.com", true},
		{"ideas-are-more-powerful-than-guns", true},
		{"testbucket", true},
		{"1bucket", true},
		{"bucket1", true},
		// cases for which test should fail.
		// passing invalid bucket names.
		{"------", false},
		{"my..bucket", false},
		{"192.168.1.1", false},
		{"$this-is-not-valid-too", false},
		{"contains-$-dollar", false},
		{"contains-^-carret", false},
		{"contains-$-dollar", false},
		{"contains-$-dollar", false},
		{"......", false},
		{"", false},
		{"a", false},
		{"ab", false},
		{".starts-with-a-dot", false},
		{"ends-with-a-dot.", false},
		{"ends-with-a-dash-", false},
		{"-starts-with-a-dash", false},
		{"THIS-BEINGS-WITH-UPPERCASe", false},
		{"tHIS-ENDS-WITH-UPPERCASE", false},
		{"ThisBeginsAndEndsWithUpperCase", false},
		{"una ñina", false},
		{"lalalallalallalalalallalallalala-theString-size-is-greater-than-64", false},
	}

	for i, testCase := range testCases {
		isValidBucketName := IsValidBucketName(testCase.bucketName)
		if testCase.shouldPass && !isValidBucketName {
			t.Errorf("Test case %d: Expected \"%s\" to be a valid bucket name", i+1, testCase.bucketName)
		}
		if !testCase.shouldPass && isValidBucketName {
			t.Errorf("Test case %d: Expected bucket name \"%s\" to be invalid", i+1, testCase.bucketName)
		}
	}
}

// Tests for validate object name.
func TestIsValidObjectName(t *testing.T) {
	testCases := []struct {
		objectName string
		shouldPass bool
	}{
		// cases which should pass the test.
		// passing in valid object name.
		{"object", true},
		{"The Shining Script <v1>.pdf", true},
		{"Cost Benefit Analysis (2009-2010).pptx", true},
		{"117Gn8rfHL2ACARPAhaFd0AGzic9pUbIA/5OCn5A", true},
		{"SHØRT", true},
		{"f*le", true},
		{"contains-^-carret", true},
		{"contains-|-pipe", true},
		{"contains-\"-quote", true},
		{"contains-`-tick", true},
		{"There are far too many object names, and far too few bucket names!", true},
		// cases for which test should fail.
		// passing invalid object names.
		{"", false},
		{"a/b/c/", false},
		{"/a/b/c", false},
		{"contains-\\-backslash", false},
		{string([]byte{0xff, 0xfe, 0xfd}), false},
	}

	for i, testCase := range testCases {
		isValidObjectName := IsValidObjectName(testCase.objectName)
		if testCase.shouldPass && !isValidObjectName {
			t.Errorf("Test case %d: Expected \"%s\" to be a valid object name", i+1, testCase.objectName)
		}
		if !testCase.shouldPass && isValidObjectName {
			t.Errorf("Test case %d: Expected object name \"%s\" to be invalid", i+1, testCase.objectName)
		}
	}
}

// Tests rangeReader.
func TestRangeReader(t *testing.T) {
	testCases := []struct {
		data   string
		minLen int64
		maxLen int64
		err    error
	}{
		{"1234567890", 0, 15, nil},
		{"1234567890", 0, 10, nil},
		{"1234567890", 0, 5, toObjectErr(errDataTooLarge, "test", "test")},
		{"123", 5, 10, toObjectErr(errDataTooSmall, "test", "test")},
		{"123", 2, 10, nil},
	}

	for i, test := range testCases {
		r := strings.NewReader(test.data)
		_, err := ioutil.ReadAll(&rangeReader{
			Reader: r,
			Min:    test.minLen,
			Max:    test.maxLen,
		})
		if toObjectErr(err, "test", "test") != test.err {
			t.Fatalf("test %d failed: expected %v, got %v", i+1, test.err, err)
		}
	}
}

// Tests getCompleteMultipartMD5
func TestGetCompleteMultipartMD5(t *testing.T) {
	testCases := []struct {
		parts          []completePart
		expectedResult string
		expectedErr    string
	}{
		// Wrong MD5 hash string
		{[]completePart{{ETag: "wrong-md5-hash-string"}}, "", "encoding/hex: odd length hex string"},

		// Single completePart with valid MD5 hash string.
		{[]completePart{{ETag: "cf1f738a5924e645913c984e0fe3d708"}}, "10dc1617fbcf0bd0858048cb96e6bd77-1", ""},

		// Multiple completePart with valid MD5 hash string.
		{[]completePart{{ETag: "cf1f738a5924e645913c984e0fe3d708"}, {ETag: "9ccbc9a80eee7fb6fdd22441db2aedbd"}}, "0239a86b5266bb624f0ac60ba2aed6c8-2", ""},
	}

	for i, test := range testCases {
		result, err := getCompleteMultipartMD5(test.parts)
		if result != test.expectedResult {
			t.Fatalf("test %d failed: expected: result=%v, got=%v", i+1, test.expectedResult, result)
		}
		errString := ""
		if err != nil {
			errString = err.Error()
		}
		if errString != test.expectedErr {
			t.Fatalf("test %d failed: expected: err=%v, got=%v", i+1, test.expectedErr, err)
		}
	}
}
