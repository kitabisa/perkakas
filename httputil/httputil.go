package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kitabisa/perkakas/v2/random"
)

const passwordPattern = `(\\{0,1}"password\\{0,1}"):\s*\\{0,1}"(.*?)\\{0,1}"`

var passRemover *regexp.Regexp

func init() {
	passRemover = regexp.MustCompile(passwordPattern)
}

func ReadRequestBody(req *http.Request) (bodyString string) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	bodyString = string(bodyBytes)
	return
}

// ExcludeSensitiveHeader exclude sensitive header. Currently, sensitive header only Authorization
func ExcludeSensitiveHeader(header http.Header) (h http.Header) {
	h = make(http.Header)
	for k, v := range header {
		h[k] = v
	}

	h.Del("Authorization")
	return
}

func ExcludeSensitiveRequestBody(body *string) {
	result := passRemover.ReplaceAllString(*body, "")
	*body = result
}

func KitabisaHeader(req *http.Request, clientName, clientVersion, requestID string) *http.Request {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	if !govalidator.IsSemver(clientVersion) {
		clientVersion = "1.0.0"
	}

	if requestID == "" {
		var err error
		requestID, err = random.UUID()
		if err != nil {
			requestID = fmt.Sprintf("%s-%s-error-generate-uuid", clientName, clientVersion)
		}
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	req.Header.Set("X-Ktbs-Request-ID", requestID)
	req.Header.Set("X-Ktbs-Client-Name", clientName)
	req.Header.Set("X-Ktbs-Client-Version", clientVersion)
	req.Header.Set("X-Ktbs-Time", timestamp)
	return req
}

func IsSuccess(code int) bool {
	return code >= 200 && code <= 299
}

func IsClientError(code int) bool {
	return code >= 400 && code <= 499
}

func IsRedirection(code int) bool {
	return code >= 300 && code <= 399
}

func IsServerError(code int) bool {
	return code >= 500 && code <= 599
}
