package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func clean(bb []byte, u string) string {
	s := strings.TrimSpace(string(bb))
	s = strings.TrimPrefix(s, "error: ")
	s = strings.TrimSuffix(s, fmt.Sprintf(": %s", u))
	return s
}

func streamlink(u string) (string, error) {
	cmd := exec.Command(
		"/usr/bin/streamlink",
		"--stream-url",
		u,
	)
	out, err := cmd.Output()
	s := clean(out, u)
	if err != nil {
		return "", errors.New(s)
	}
	return s, nil
}
