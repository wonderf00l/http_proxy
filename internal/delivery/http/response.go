package http

import (
	"encoding/json"
	"errors"
	"net/http"

	errPkg "github.com/wonderf00l/http_proxy/internal/errors"
)

type JSONErrResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"message"`
}

func GetCodeStatusHttp(err error) (ErrCode string, httpStatus int) {
	var declaredErr errPkg.DeclaredError
	if errors.As(err, &declaredErr) {
		switch declaredErr.Type() {
		case errPkg.ErrInvalidInput:
			return "bad_input", http.StatusBadRequest
		case errPkg.ErrNotFound:
			return "not_found", http.StatusNotFound
		}
	}

	return "internal_error", http.StatusInternalServerError
}

func (h *HandlerHTTP) responseErr(w http.ResponseWriter, r *http.Request, err error) {
	code, status := GetCodeStatusHttp(err)
	var msg string
	if status == http.StatusInternalServerError {
		h.logger.Warnf("unexpected application error: %s", err.Error())
		msg = "internal error occured"
	} else {
		msg = err.Error()
	}

	res := JSONErrResponse{
		Status: "error",
		Msg:    msg,
		Code:   code,
	}
	resBytes, _ := json.Marshal(res)

	w.WriteHeader(status)
	w.Write(resBytes)
}
