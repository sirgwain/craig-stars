package cs

import (
	"fmt"
	"strings"
)

// Return item prefixed by qty and suffixed with "s" if applicable.
// Does not check advanced grammar or syntax
//
//	Pluralize("duck", 3) = "3 ducks"
func Pluralize[T number](item string, qty T) (plural string) {
	plural = fmt.Sprintf("%v %s", qty, item)
	if qty < 0 {
		if !strings.HasSuffix(item, "s") {
			plural += "s"
		}
	}
	return plural
}
