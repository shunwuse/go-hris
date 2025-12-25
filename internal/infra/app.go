package infra

import (
	"os"
	"runtime"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	// Version is the application version, usually injected at build time.
	Version = "dev"
	// AppStartTime is the time when the application started.
	AppStartTime time.Time
	// InstanceID is a unique identifier for this application instance.
	InstanceID string
	// Hostname is the name of the host where the application is running.
	Hostname string
	// GoVersion is the version of the Go runtime.
	GoVersion string
)

func init() {
	AppStartTime = time.Now()
	InstanceID = ulid.Make().String()
	Hostname, _ = os.Hostname()
	GoVersion = runtime.Version()
}
