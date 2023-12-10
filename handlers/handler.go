package handler

import (
	"fmt"
	"go-nutritioncalculator2/errs"
	"net/http"
)

func handlerError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case errs.AppError:
		w.WriteHeader(e.Code)
		fmt.Fprintln(w, e)
	}
}
