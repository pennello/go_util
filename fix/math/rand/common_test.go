// chris 072815 Common testing code.  Seed random number generator with
// current time.

package rand

import (
	"os"
	"testing"
	"time"

	"math/rand"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}
