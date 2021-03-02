package params

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const (
	defRate     = 1
	defInflight = 1
)

var (
	ErrNoReplacementMarker = errors.New("no replacement marker")
)

// Get get application params
func Get(args []string) (int, int, string, error) {
	rate := flag.Int("rate", defRate, fmt.Sprintf("Rate limit (default %d)", defRate))
	inflight := flag.Int("inflight", defInflight, fmt.Sprintf("Inflight (default %d)", defInflight))
	flag.Parse()

	cmd := strings.Join(flag.Args(), " ")
	fcmd := strings.Replace(cmd, "{}", "%s", 1)
	if fcmd == cmd {
		return 0, 0, "", ErrNoReplacementMarker
	}

	return *rate, *inflight, fcmd, nil
}
