package request_utils

import (
	"net/http"
	"strconv"
)

type RequestUtils struct{}

func NewRequestUtils() *RequestUtils {
	return &RequestUtils{}
}

func (u *RequestUtils) GetQueryInt(r *http.Request, param string, def int) (int, error) {
	paramStr := r.URL.Query().Get(param)
	if paramStr == "" {
		return def, nil
	}
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	return paramInt, nil
}

func (u *RequestUtils) GetContextInt(r *http.Request, key any, def int) int {
	i, ok := r.Context().Value(key).(int)
	if !ok {
		return def
	}
	return i
}
