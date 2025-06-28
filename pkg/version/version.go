package version

import (
	"fmt"
	"runtime/debug"
)

var Version string
var Revision string
var CommitTime string
var Dirty string

func GetVersion() string {
	return Version
}

func GetRevision() string {
	return Revision
}

func GetVersionRevision() string {
	return fmt.Sprintf("%s-%s", Version, Revision)
}

type VersionInfo struct {
	Version  string `json:"version"`
	Revision string `json:"revision"`
}

func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		Version:  Version,
		Revision: Revision,
	}
}

type BuildInfo struct {
	Module      debug.Module `json:"module"`
	VCS         string       `json:"vcs"`
	VCSRevision string       `json:"vcs_revision"`
	VCSTime     string       `json:"vcs_time"`
	VCSModified string       `json:"vcs_modified"`
	Tags        []string     `json:"tags"`
}

func GetBuildInfo() *BuildInfo {
	// Leaving this here for posterity

	// buildInfo, ok := debug.ReadBuildInfo()
	// if !ok {
	// 	slog.Global().Error("failed to read build info")
	// 	return &BuildInfo{}
	// }

	// vcs, _ := lo.Find(buildInfo.Settings, func(item debug.BuildSetting) bool {
	// 	return strings.EqualFold(item.Key, "vcs")
	// })
	// vcsRevision, _ := lo.Find(buildInfo.Settings, func(item debug.BuildSetting) bool {
	// 	return strings.EqualFold(item.Key, "vcs.revision")
	// })
	// vcsTime, _ := lo.Find(buildInfo.Settings, func(item debug.BuildSetting) bool {
	// 	return strings.EqualFold(item.Key, "vcs.time")
	// })
	// vcsModified, _ := lo.Find(buildInfo.Settings, func(item debug.BuildSetting) bool {
	// 	return strings.EqualFold(item.Key, "vcs.modified")
	// })
	// tags, _ := lo.Find(buildInfo.Settings, func(item debug.BuildSetting) bool {
	// 	return strings.EqualFold(item.Key, "tags") || strings.EqualFold(item.Key, "-tags")
	// })

	return &BuildInfo{
		Module:      debug.Module{},
		VCS:         "git",
		VCSRevision: Revision,
		VCSTime:     CommitTime,
		VCSModified: Dirty,
		Tags:        []string{Version},
	}
}
