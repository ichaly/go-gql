package transport

import (
	"encoding/json"
	"github.com/ichaly/go-gql"
	"io"
)

func readJson(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func writeJson(w io.Writer, response *graphql.Response) {
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(b)
}
