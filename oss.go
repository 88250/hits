// Hits - GitHub repository hits counter.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"os"
	"os/exec"
	u "os/user"
	"runtime"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}

func IsWindows() bool {
	return "windows" == runtime.GOOS
}

func UserHome() string {
	user, err := u.Current()
	if nil == err {
		return user.HomeDir
	}

	if IsWindows() {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() string {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		logger.Errorf("get user home path failed [%s]", err)

		return ""
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		logger.Errorf("blank output when reading home directory")

		return ""
	}

	return result
}

func homeWindows() string {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		logger.Errorf("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")

		return ""
	}

	return home
}
