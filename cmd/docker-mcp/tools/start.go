package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/docker/mcp-gateway/pkg/desktop"
	"github.com/docker/mcp-gateway/pkg/logs"
)

func start(ctx context.Context, version string, gatewayArgs []string, verbose bool) (*mcp.ClientSession, error) {
	var args []string
	if version == "2" {
		if verbose {
			args = []string{"mcp", "gateway", "run", "--verbose"}
		} else {
			args = []string{"mcp", "gateway", "run"}
		}
	} else {
		// Usa detecção automática para determinar o endereço do host
		hostAddress := desktop.GetHostAddress()
		socatTarget := fmt.Sprintf("TCP:%s:8811", hostAddress)
		args = []string{"run", "-i", "--rm", "alpine/socat", "STDIO", socatTarget}

		if verbose {
			prefixer := logs.NewPrefixer(os.Stderr, "- mcp-gateway: ")
			fmt.Fprintf(prefixer, "Usando endereço de host: %s\n", hostAddress)
		}
	}
	args = append(args, gatewayArgs...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	if verbose {
		cmd.Stderr = logs.NewPrefixer(os.Stderr, "- mcp-gateway: ")
	}

	c := mcp.NewClient(&mcp.Implementation{Name: "mcp-gateway-client", Version: "1.0.0"}, nil)
	transport := &mcp.CommandTransport{Command: cmd}
	session, err := c.Connect(ctx, transport, nil)
	if err != nil {
		return nil, err
	}

	return session, nil
}
