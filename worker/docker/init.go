package docker

import "github.com/lab23/abstruse/worker/config"

var (
	cfg *config.Registry
)

// Init initializes global variables
func Init(config *config.Registry) {
	cfg = config
}
