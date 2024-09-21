package generators

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateUniqueID creates a unique ID with unique information
func GenerateUniqueID(id, name string) string {

	n := shuffler(fmt.Sprintf("%s%s", strings.ReplaceAll(name, " ", ""), id))
	t := time.Now()
	tt := shuffler(fmt.Sprint(t.Unix()))

	return strings.ReplaceAll(fmt.Sprintf("%s-%s-%d", tt, n, t.Unix()), " ", "")
}

// shuffler shuffles all the words
func shuffler(s string) string {
	r := []rune(s)
	rand.Shuffle(len(r), func(i, j int) {
		r[i], r[j] = r[j], r[i]
	})
	return string(r)
}
