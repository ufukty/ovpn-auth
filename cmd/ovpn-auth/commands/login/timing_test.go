package login

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ufukty/ovpn-auth/internal/files"
	"github.com/ufukty/ovpn-auth/internal/utils"
)

func TestTimings(t *testing.T) {
	tcs := map[string][]time.Duration{
		"invalid-totp.yml":          nil,
		"invalid-username.yml":      nil,
		"invalid-password-totp.yml": nil,
	}

	db, err := files.LoadDatabase("testdata/database.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, loading database: %w", err))
	}

	for tn := range tcs {
		t.Run(tn, func(t *testing.T) {
			rqs, err := loadTestRequestsFile(filepath.Join("testdata", tn))
			if err != nil {
				t.Fatal(fmt.Errorf("loading test requests file %q: %w", tn, err))
			}
			tcs[tn] = []time.Duration{}
			for _, rq := range rqs {
				start := time.Now()
				login(db, &rq)
				tcs[tn] = append(tcs[tn], time.Since(start))
			}
		})
	}

	fmt.Println("total durations for all requests in same test set:")
	for tn, tc := range tcs {
		fmt.Printf("    %-25s => %s\n", tn, utils.Sum(tc))
	}

	fmt.Println("average durations per request:")
	for tn, tc := range tcs {
		fmt.Printf("    %-25s => %.0fms\n", tn, utils.Avg(tc)/1000000)
	}

	fmt.Println("standard deviations (amongst all requests in one set, individually):")
	for tn, tc := range tcs {
		fmt.Printf("    %-25s => %.2fms\n", tn, utils.StandardDeviation(tc)/1000000)
	}
}
