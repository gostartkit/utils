package utils

import (
	"os/exec"
	"runtime"
)

// OpenBrowser attempts to open the specified URL in the default browser of the user.
func OpenBrowser(url string) error {

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	default: // Unix-like systems: "linux", "freebsd", "openbsd", "netbsd"
		return exec.Command("xdg-open", url).Run()
	}
}
