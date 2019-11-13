package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/svipul/tinfoil-wizard/errorutil"
)

const (
	jsonPrefix = ""
	jsonIndent = "  "
)

func ErrorResponse(w http.ResponseWriter, err error) {
	code := errorutil.Code(err)
	if code == nil {
		*code = 500
	}
	http.Error(w, err.Error(), *code)
}

func JSONResponse(w http.ResponseWriter, successCode int, respPayload interface{}) {
	js, err := json.MarshalIndent(respPayload, jsonPrefix, jsonIndent)
	if err != nil {
		fmt.Printf("JSON parse error: %v", err)
		http.Error(w, "internal error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successCode)
	if js != nil {
		if _, err := w.Write(js); err != nil {
			http.Error(w, "couldn't write response", 500)
			return
		}
	}
}
