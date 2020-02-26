package system

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// The following fields are populated at build time using -ldflags -X.
// Note that DATE is omitted for reproducible builds
var (
	buildName        = "unknown"
	buildVersion     = "unknown"
	buildGitRevision = "unknown"
	buildUser        = "unknown"
	buildHostName    = "unknown"
	buildStatus      = "unknown"
	buildTime        = "unknown"
)

// BuildInfo describes version information about the binary build.
type buildInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	GitRevision   string `json:"gitRevision"`
	User          string `json:"user"`
	HostName      string `json:"hostName"`
	GolangVersion string `json:"golangVersion"`
	BuildStatus   string `json:"buildStatus"`
	BuildTime     string `json:"buildTime"`
	BuildDebug    bool   `json:"buildDebug"`
}

// RunInfo describes version information about the binary run.
type runInfo struct {
	Pid       int       `json:"pid"`
	HostName  string    `json:"hostName"`
	StartTime time.Time `json:"startTime"`
}

var (
	// Info exports the build version information.
	BuildInfo buildInfo
	RunInfo   runInfo
)

func InitRunInfo() {
	var err error
	RunInfo.Pid = os.Getpid()
	RunInfo.HostName, err = os.Hostname()
	RunInfo.StartTime = time.Now()
	if err != nil {
		// todo log
	}
}

// String produces a single-line version info
//
// This looks like:
//
// ```
// name-user@host-<version>-<git revision>-<build status>
// ```
func (b buildInfo) String() string {
	return fmt.Sprintf("%v@%v-%v-%v-%v-%v",
		b.User,
		b.HostName,
		b.Version,
		b.GitRevision,
		b.BuildStatus,
		b.BuildTime)
}

// LongForm returns a multi-line version information
//
// This looks like:
//
// ```
// Version: <version>
// GitRevision: <git revision>
// User: user@host
// GolangVersion: go1.10.2
// BuildStatus: <build status>
// ```
func (b buildInfo) LongForm() string {
	return fmt.Sprintf(`Name: %v
Version: %v
GitRevision: %v
User: %v@%v
GolangVersion: %v
BuildStatus: %v
BuildTime: %v
`,
		b.Name,
		b.Version,
		b.GitRevision,
		b.User,
		b.HostName,
		b.GolangVersion,
		b.BuildStatus,
		b.BuildTime)
}

func (r runInfo) String() string {
	return fmt.Sprintf("%v@%v-%v",
		r.Pid,
		r.HostName,
		r.StartTime,
	)
}

func (r runInfo) LongForm() string {
	return fmt.Sprintf(`Pid: %v
HostName: %v
StartTime: %v
`,
		r.Pid,
		r.HostName,
		r.StartTime,
	)
}

func init() {
	BuildInfo = buildInfo{
		Name:          buildName,
		Version:       buildVersion,
		GitRevision:   buildGitRevision,
		User:          buildUser,
		HostName:      buildHostName,
		GolangVersion: runtime.Version(),
		BuildStatus:   buildStatus,
		BuildTime:     buildTime,
	}

}
