package trade

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, r *http.Request, statusCode int, body interface{}) {
	out, err := json.MarshalIndent(body, "", "	")
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
		return
	}

	w.WriteHeader(statusCode)
	w.Write(out)
}

type errorResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"-"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}

func RenderError(w http.ResponseWriter, r *http.Request, code int, svrErr error, format string, args ...any) {
	var err error
	errResp := newErrorResponse(code, svrErr, format, args...)
	respBody, err := json.Marshal(errResp)
	if err != nil {
		mustWriteError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
	return
}

func newErrorResponse(code int, err error, format string, args ...any) *errorResponse {
	return &errorResponse{
		Err:            err,
		HTTPStatusCode: code,
		ErrorText:      fmt.Sprintf(format, args...),
	}
}
func mustWriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("error encountered while attempting to write error: " + err.Error()))
	return
}
