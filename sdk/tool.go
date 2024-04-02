package sdk

import (
	"os"
	"os/exec"
)

// Clear 清屏
func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
