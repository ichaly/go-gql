package introspection

import (
	"fmt"
	"github.com/ichaly/go-gql/internal/introspection/location"
	"testing"
)

func TestDirectiveLocation(t *testing.T) {
	fmt.Printf("%v", location.QUERY)
}
