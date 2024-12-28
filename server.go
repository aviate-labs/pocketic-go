package pocketic

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type printWriter struct{}

func (w printWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

type server struct {
	port int
	cmd  *exec.Cmd
}

func newServer(opts ...serverOption) (*server, error) {
	config := serverConfig{}
	for _, fn := range opts {
		fn(&config)
	}

	// Try to find the pocket-ic binary.
	binPath, err := exec.LookPath("pocket-ic-server")
	if err != nil {
		if binPath, err = exec.LookPath("pocket-ic"); err != nil {
			// If the binary is not found, try to find it in the POCKET_IC_BIN environment variable.
			if pathEnv := os.Getenv("POCKET_IC_BIN"); pathEnv != "" {
				binPath = pathEnv
			} else {
				binPath = "./pocket-ic"
				if _, err := os.Stat(binPath); err != nil {
					return nil, fmt.Errorf("pocket-ic binary not found: %v", err)
				}
			}
		}
	}

	versionCmd := exec.Command(binPath, "--version")
	rawVersion, err := versionCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get pocket-ic version: %v", err)
	}
	version := strings.TrimPrefix(strings.TrimSpace(string(rawVersion)), "pocket-ic-server ")
	if !strings.HasPrefix(version, "4.") {
		return nil, fmt.Errorf("unsupported pocket-ic version, must be v4.x: %s", version)
	}

	pid := os.Getpid()
	cmdArgs := []string{"--pid", strconv.Itoa(pid)}
	if config.ttl != nil {
		cmdArgs = append(cmdArgs, "--ttl", strconv.Itoa(*config.ttl))
	}
	cmd := exec.Command(binPath, cmdArgs...)
	if os.Getenv("POCKET_IC_MUTE_SERVER") == "" {
		cmd.Stdout = new(printWriter)
		cmd.Stderr = new(printWriter)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start pocket-ic: %v", err)
	}

	tmpDir := os.TempDir()
	readyFile := path.Join(tmpDir, fmt.Sprintf("pocket_ic_%d.ready", pid))
	stopAt := time.Now().Add(5 * time.Second)
	for _, err := os.Stat(readyFile); os.IsNotExist(err); _, err = os.Stat(readyFile) {
		time.Sleep(100 * time.Millisecond)
		if time.Now().After(stopAt) {
			return nil, fmt.Errorf("pocket-ic did not start in time, %s not found", readyFile)
		}
	}

	portFile := path.Join(tmpDir, fmt.Sprintf("pocket_ic_%d.port", pid))
	f, err := os.OpenFile(portFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open port file: %v", err)
	}
	rawPort, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read port file: %v", err)
	}
	port, err := strconv.Atoi(string(rawPort))
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %v", err)
	}

	return &server{
		port: port,
		cmd:  cmd,
	}, nil
}

func (s server) URL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", s.port)
}

type serverConfig struct {
	ttl *int
}

type serverOption func(*serverConfig)

// withTTL sets the time-to-live for the pocket-ic server, in seconds.
func withTTL(ttl int) serverOption {
	return func(c *serverConfig) {
		c.ttl = &ttl
	}
}
