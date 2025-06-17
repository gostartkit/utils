package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// OpenBrowser attempts to open the specified URL in the default browser of the user.
func OpenBrowser(uri string) error {

	if !isValidURI(uri) {
		return fmt.Errorf("invalid URI: %s", uri)
	}

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", uri).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Run()
	default: // Unix-like systems: "linux", "freebsd", "openbsd", "netbsd"
		if err := exec.Command("xdg-open", uri).Run(); err != nil {
			// wsl
			uri = escapeForCmd(uri)
			if err2 := exec.Command("cmd.exe", "/C", "start", "chrome", uri).Run(); err2 != nil {
				if err3 := exec.Command("cmd.exe", "/C", "start", uri).Run(); err3 != nil {
					return err
				}
			}
		}
		return nil
	}
}

func escapeForCmd(uri string) string {
	// Escape CMD special characters: &, |, >, <
	replacements := map[string]string{
		"&": "^&",
		"|": "^|",
		">": "^>",
		"<": "^<",
	}
	result := uri
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

func isValidURI(uri string) bool {
	return strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://")
}
