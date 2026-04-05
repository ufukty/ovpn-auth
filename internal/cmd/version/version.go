package version

import (
	"fmt"

	"github.com/ufukty/ovpn-auth/internal/version"
)

func Run() error {
	v, err := version.OfBuild()
	if err != nil {
		return fmt.Errorf("digging build details: %w", err)
	}
	fmt.Println(v)
	return nil
}
