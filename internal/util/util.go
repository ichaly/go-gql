package util

import (
	"encoding/json"
	"fmt"
	"github.com/ichaly/go-gql/internal/executor"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"io"
	"net/http"
	"time"
)

var Now = time.Now

func ReadJson(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func WriteJson(w io.Writer, res *executor.GqlResult) {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}

func SendErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendError(w, code, &gqlerror.Error{Message: fmt.Sprintf(format, args...)})
}

func SendError(w http.ResponseWriter, code int, errors ...*gqlerror.Error) {
	w.WriteHeader(code)
	b, err := json.Marshal(&executor.GqlResult{Errors: errors})
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
