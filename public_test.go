package strconf_test

import "github.com/goappconf/strconf"

// Verify that consumers can access the API without running remote scripts.
var _ func() error = strconf.Run
