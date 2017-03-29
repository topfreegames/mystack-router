package extensions

import (
	"fmt"
	"os/exec"
	"strings"
)

func getNginxBinary() (string, error) {
	out, err := exec.Command("which", "nginx").Output()
	if err != nil {
		return "", err
	}

	nginxBinary := string(out)
	if strings.Contains(nginxBinary, "not found") {
		return "", fmt.Errorf("nginx is not installed")
	}

	return nginxBinary, nil
}

//Start nginx
func Start() error {
	return nil
}
