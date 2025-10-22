package main

import (
	"fmt"
	"runtime/debug"
)

func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		// Return the module version if available (from git tag)
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}

		// Development build - show commit info
		var revision, modified string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value[:7]
			case "vcs.modified":
				if setting.Value == "true" {
					modified = "-dirty"
				}
			}
		}

		if revision != "" {
			return fmt.Sprintf("dev-%s%s", revision, modified)
		}
	}

	return "dev"
}
