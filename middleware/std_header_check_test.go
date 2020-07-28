package middleware

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/signature"
	"github.com/kitabisa/perkakas/v2/structs"
	uuid "github.com/satori/go.uuid"
)

var testHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, client")
})

func TestStdHeaderValidation(t *testing.T) {
	hctx := phttp.NewContextHandler(structs.Meta{})
	headerCheck := NewHeaderCheck(hctx, "key")
	ts := httptest.NewServer(headerCheck(testHandler))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.FailNow()
	}

	theTime := time.Now().UTC().Unix()

	sign := signature.GenerateHmac(fmt.Sprintf("%s%d", "kitabisa-apps", theTime), "key")

	req.Header.Add("X-Ktbs-Request-ID", uuid.NewV4().String())
	req.Header.Add("X-Ktbs-Api-Version", "1.0.1")
	req.Header.Add("X-Ktbs-Client-Version", "1.1.1")
	req.Header.Add("X-Ktbs-Platform-Name", "android")
	req.Header.Add("X-Ktbs-Client-Name", "kitabisa-apps")
	req.Header.Add("X-Ktbs-Signature", sign)
	req.Header.Add("X-Ktbs-Time", strconv.FormatInt(theTime, 10))
	req.Header.Add("Authorization", "Bearer")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.FailNow()
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal(res.StatusCode, string(greeting))
	}
}
