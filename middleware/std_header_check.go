package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/signature"
	"github.com/kitabisa/perkakas/v2/structs"
)

type Header struct {
	XKtbsRequestID     string `valid:"uuidv4,required"`
	XKtbsAPIVersion    string `valid:"semver,required"`
	XKtbsClientVersion string `valid:"semver,required"`
	XKtbsPlatformName  string `valid:"required"`
	XKtbsClientName    string `valid:"required"`

	// Optional
	XKtbsSignature string `valid:"optional"`
	XKtbsTime      string `valid:"int,optional"`
	Authorization  string `valid:"optional"`
}

func NewHeaderCheck(hctx phttp.HttpHandlerContext, secretKey string) func(next http.Handler) http.Handler {
	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := Header{
				XKtbsRequestID:     r.Header.Get("X-Ktbs-Request-ID"),
				XKtbsAPIVersion:    r.Header.Get("X-Ktbs-Api-Version"),
				XKtbsClientVersion: r.Header.Get("X-Ktbs-Client-Version"),
				XKtbsPlatformName:  r.Header.Get("X-Ktbs-Platform-Name"),
				XKtbsClientName:    r.Header.Get("X-Ktbs-Client-Name"),
				XKtbsSignature:     r.Header.Get("X-Ktbs-Signature"),
				XKtbsTime:          r.Header.Get("X-Ktbs-Time"),
				Authorization:      r.Header.Get("Authorization"),
			}

			_, err := govalidator.ValidateStruct(header)
			if err != nil {
				writer.WriteError(w, structs.ErrInvalidHeader)
				return
			}

			theTime, err := strconv.ParseInt(header.XKtbsTime, 10, 64)
			if err != nil {
				writer.WriteError(w, structs.ErrInvalidHeader)
				return
			}

			// delay request should not be more than 1 min
			if theTime+60 < time.Now().UTC().Unix() {
				writer.WriteError(w, structs.ErrInvalidHeaderTime)
				return
			}

			data := fmt.Sprintf("%s%s", header.XKtbsClientName, header.XKtbsTime)
			match := signature.IsMatchHmac(data, header.XKtbsSignature, secretKey)
			if !match {
				writer.WriteError(w, structs.ErrInvalidHeaderSignature)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
