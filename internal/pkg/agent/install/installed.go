// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package install

import (
	"os"
	"path/filepath"

	"github.com/kardianos/service"

	"github.com/elastic/elastic-agent/internal/pkg/agent/application/paths"
	"github.com/elastic/elastic-agent/internal/pkg/agent/install/pkgmgr"
)

// StatusType is the return status types.
type StatusType int

const (
	// NotInstalled returned when Elastic Agent is not installed.
	NotInstalled StatusType = iota
	// Installed returned when Elastic Agent is installed correctly.
	Installed
	// Broken returned when Elastic Agent is installed but broken.
	Broken
	// PackageInstall returned when the Elastic agent has been installed through a package manager (deb/rpm)
	PackageInstall
)

// Status returns the installation status of Agent.
func Status(topPath string) (StatusType, string) {
	expected := filepath.Join(topPath, paths.BinaryName)
	status, reason := checkService(topPath)
	if pkgmgr.InstalledViaExternalPkgMgr() {
		if status == Installed {
			return PackageInstall, "service running"
		}
		return PackageInstall, "service not running"
	}
	_, err := os.Stat(expected)
	if os.IsNotExist(err) {
		if status == Installed {
			// service installed, but no install path
			return Broken, "service exists but installation path is missing"
		}
		return NotInstalled, "no install path or service"
	}
	if status == NotInstalled {
		// install path present, but not service
		return Broken, reason
	}
	return Installed, ""
}

// checkService only checks the status of the service.
func checkService(topPath string) (StatusType, string) {
	svc, err := newService(topPath)
	if err != nil {
		return NotInstalled, "unable to check service status"
	}
	status, _ := svc.Status()
	if status == service.StatusUnknown {
		return NotInstalled, "service is not installed"
	}
	return Installed, ""
}
