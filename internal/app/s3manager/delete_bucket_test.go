package s3manager_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mastertinner/s3manager/internal/app/s3manager"
	"github.com/mastertinner/s3manager/internal/app/s3manager/mocks"
	"github.com/matryer/is"
)

func TestHandleDeleteBucket(t *testing.T) {
	cases := map[string]struct {
		removeBucketFunc     func(string) error
		expectedStatusCode   int
		expectedBodyContains string
	}{
		"deletes an existing bucket": {
			removeBucketFunc: func(string) error {
				return nil
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedBodyContains: "",
		},
		"returns error if there is an S3 error": {
			removeBucketFunc: func(string) error {
				return errors.New("mocked S3 error")
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: http.StatusText(http.StatusInternalServerError),
		},
	}

	for tcID, tc := range cases {
		t.Run(tcID, func(t *testing.T) {
			is := is.New(t)

			s3 := &mocks.S3Mock{
				RemoveBucketFunc: tc.removeBucketFunc,
			}

			req, err := http.NewRequest(http.MethodDelete, "/api/buckets/bucketName", nil)
			is.NoErr(err)

			rr := httptest.NewRecorder()
			handler := s3manager.HandleDeleteBucket(s3)

			handler.ServeHTTP(rr, req)
			resp := rr.Result()

			is.Equal(tc.expectedStatusCode, resp.StatusCode)                     // status code
			is.True(strings.Contains(rr.Body.String(), tc.expectedBodyContains)) // body
		})
	}
}