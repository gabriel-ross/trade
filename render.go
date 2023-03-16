package trade

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RenderService struct{}

func (rs *RenderService) RenderJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, body interface{}) {
	out, err := json.MarshalIndent(body, "", "	")
	if err != nil {
		rs.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
		return
	}

	w.WriteHeader(httpStatusCode)
	w.Write(out)
}

func (rs *RenderService) RenderError(w http.ResponseWriter, r *http.Request, svrErr error, code int, format string, args ...any) {
	var err error
	errResp := rs.newErrorResponse(code, svrErr, format, args...)
	respBody, err := json.Marshal(errResp)
	if err != nil {
		rs.mustWriteError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
	return
}

func (rs *RenderService) newErrorResponse(code int, err error, format string, args ...any) *errorResponse {
	return &errorResponse{
		Err:            err,
		HTTPStatusCode: code,
		ErrorText:      fmt.Sprintf(format, args...),
	}
}
func (rs *RenderService) mustWriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("error encountered while attempting to write error: " + err.Error()))
	return
}

type errorResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"-"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}
