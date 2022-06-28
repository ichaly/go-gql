package util

import (
	"encoding/json"
	"github.com/ichaly/go-gql/types"
	"io"
	"time"
)

var Now = time.Now

func ReadJson(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func WriteJson(w io.Writer, response *types.GqlResult) {
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
