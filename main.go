package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func fail(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func shouldFetchKubeconfig(kubeconfig string) bool {
	stat, err := os.Stat(kubeconfig)
	// Fetch kubeconfig if doesn't exist.
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	// Renew kubeconfig if more than a week old.
	if time.Since(stat.ModTime()) > 7*24*time.Hour {
		return true
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fail("Usage: %s vsn\n", os.Args[0])
	}

	vsn := os.Args[1]

	home, err := os.UserHomeDir()
	if err != nil {
		fail("Could not get home directory.\n")
	}

	kubeconfig := filepath.Join(home, ".kube", fmt.Sprintf("k3s-%s.yaml", vsn))

	if shouldFetchKubeconfig(kubeconfig) {
		fmt.Printf("Fetching node's kubeconfig...\n")
		cmd := exec.Command("scp", fmt.Sprintf("node-%s:/etc/rancher/k3s/k3s.yaml", vsn), kubeconfig)
		if err := cmd.Run(); err != nil {
			fail("Could not fetch node's kubeconfig!\n")
		}
	}

	fmt.Printf("Starting tunnel...\n")
	tunnelCmd := exec.Command("ssh", fmt.Sprintf("node-%s", vsn), "-D", "1080", "-N")
	if err := tunnelCmd.Start(); err != nil {
		fail("Could not start tunnel!\n")
	}
	// TOOD Wait for this to succeed before starting shell!
	defer tunnelCmd.Process.Kill()

	env := os.Environ()
	env = append(env, fmt.Sprintf("KUBECONFIG=%s", kubeconfig))
	env = append(env, "HTTPS_PROXY=socks5://localhost:1080")

	fmt.Printf("Starting attached environment...\n")
	// TODO Make shell configurable.
	shellCmd := exec.Command("/bin/bash")
	shellCmd.Env = env
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr
	shellCmd.Stdout = os.Stdout
	shellCmd.Stdin = os.Stdin
	shellCmd.Run()
}
