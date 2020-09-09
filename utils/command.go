package utils

import (
	"io"
	"os"
	"os/exec"
	"bytes"
)

func ExecCmd(command, dir string, env []string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Dir = dir
	cmd.Env = env
	var stdBuffer bytes.Buffer
	stdWriter := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = stdWriter
	cmd.Stderr = stdWriter
	e := cmd.Run()
	return e
}
