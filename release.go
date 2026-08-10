package uax

import (
	"hash/fnv"
	"strconv"
	"sync"

	"github.com/bakhod1r/uax/internal/table"
)

// Release is the version of this library, matching its Git tag.
//
// It is a constant rather than build-time linker flags so that `go get` of a
// tagged module reports the same value however it was built.
const Release = "0.1.0"

// DatabaseVersion is a fingerprint of the compiled rule database.
//
// A classification is only reproducible if you know what it was classified
// against, and the rules change more often than the code does. Logging this
// alongside a verdict makes a decision from six months ago explicable: two
// runs that report the same fingerprint saw the same rules.
//
// It is derived from the rules themselves rather than maintained by hand, so
// it cannot be forgotten. The value is stable for a given database and carries
// no meaning beyond equality — it is not ordered, and a difference says the
// rules changed without saying how.
func DatabaseVersion() string {
	databaseVersionOnce.Do(func() {
		h := fnv.New64a()
		for _, group := range [][]table.Rule{
			table.Browsers, table.Engines, table.OSes, table.Bots,
			table.Apps, table.Automations, table.Clients,
		} {
			for i := range group {
				r := &group[i]
				// ID, the literal it matches on, and what it claims. Anything
				// that changes an answer changes the fingerprint.
				_, _ = h.Write([]byte(r.ID))
				_, _ = h.Write([]byte(r.Token))
				_, _ = h.Write([]byte(r.Prefix))
				_, _ = h.Write([]byte(r.Contains))
				_, _ = h.Write([]byte(r.Name))
				_, _ = h.Write([]byte(r.Brand))
				_, _ = h.Write([]byte(strconv.Itoa(r.Priority)))
			}
		}
		databaseVersion = strconv.FormatUint(h.Sum64(), 16)
	})
	return databaseVersion
}

var (
	databaseVersionOnce sync.Once
	databaseVersion     string
)
