package utils

import (
	"fmt"
	"math/rand"
)

func Generate6DigitCode() string {
	code := rand.Intn(900000) + 100000
	
	if IsDevMode() {
		// For development mode, always return 888888
		code = 888888
	}

	return fmt.Sprintf("%06d", code)
}

