package s3manager_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/mastertinner/s3manager"
	"github.com/stretchr/testify/assert"
)

func TestCreateBucketHandler(t *testing.T) {
	assert := assert.New(t)

	cases := map[string]struct {
		s3                   S3
		body                 string
		expectedStatusCode   int
		expectedBodyContains string
	}{
		"success": {
			s3:                   &s3Mock{},
			body:                 "{\"name\":\"myBucket\"}",
			expectedStatusCode:   http.StatusCreated,
			expectedBodyContains: "{\"name\":\"myBucket\",\"creationDate\":\"0001-01-01T00:00:00Z\"}\n",
		},
		"empty request": {
			s3:                   &s3Mock{},
			body:                 "",
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedBodyContains: http.StatusText(http.StatusUnprocessableEntity),
		},
		"malformed request": {
			s3:                   &s3Mock{},
			body:                 "}",
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedBodyContains: http.StatusText(http.StatusUnprocessableEntity),
		},
		"s3 error": {
			s3: &s3Mock{
				Err: errors.New("mocked S3 error"),
			},
			body:                 "{\"name\":\"myBucket\"}",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: http.StatusText(http.StatusInternalServerError),
		},
	}

	for tcID, tc := range cases {
		t.Run(tcID, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/buckets", bytes.NewBufferString(tc.body))
			assert.NoError(err, tcID)

			rr := httptest.NewRecorder()
			handler := CreateBucketHandler(tc.s3)

			handler.ServeHTTP(rr, req)
			resp := rr.Result()

			assert.Equal(tc.expectedStatusCode, resp.StatusCode, tcID)
			assert.Contains(rr.Body.String(), tc.expectedBodyContains, tcID)
		})
	}
}