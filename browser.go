package utils

import (
	"os/exec"
	"runtime"
)

// OpenBrowser attempts to open the specified URL in the default browser of the user.
func OpenBrowser(url string) error {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // Unix-like systems: "linux", "freebsd", "openbsd", "netbsd"
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Run()
}
