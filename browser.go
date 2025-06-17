package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

// OpenBrowser attempts to open the specified URL in the default browser of the user.
func OpenBrowser(uri string) error {

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", uri).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Run()
	default: // Unix-like systems: "linux", "freebsd", "openbsd", "netbsd"
		if err := exec.Command("xdg-open", uri).Run(); err != nil {
			// wsl
			uri = strings.Replace(uri, "&", "^&", -1)
			if err2 := exec.Command("cmd.exe", "/C", "start", "chrome", uri).Run(); err2 != nil {
				if err3 := exec.Command("cmd.exe", "/C", "start", uri).Run(); err3 != nil {
					return err
				}
			}
		}
		return nil
	}
}
