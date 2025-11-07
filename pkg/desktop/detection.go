package desktop

import (
	"os"
	"runtime"

	"github.com/docker/mcp-gateway/pkg/log"
)

// IsNativeMode detecta automaticamente se o sistema está rodando em Linux nativo (sem Docker Desktop)
// A detecção segue estas regras:
// 1. Se DOCKER_MCP_NATIVE_MODE estiver definida, respeita o valor manual (override)
// 2. Se não estiver em Linux, retorna false
// 3. Se estiver em contêiner (DOCKER_MCP_IN_CONTAINER=1), retorna false
// 4. Verifica se o socket do Docker Desktop não existe
// 5. Verifica se o Docker Engine está disponível
func IsNativeMode() bool {
	// 1. Respeita override manual via variável de ambiente
	if env := os.Getenv("DOCKER_MCP_NATIVE_MODE"); env != "" {
		return env == "1"
	}

	// 2. Apenas Linux pode ter modo nativo
	if runtime.GOOS != "linux" {
		return false
	}

	// 3. Se estiver em contêiner, não é modo nativo
	if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" {
		return false
	}

	// 4. Verifica se o socket do Docker Desktop não existe
	dockerDesktopSocketExists := false
	if _, err := os.Stat("/run/host-services/backend.sock"); err == nil {
		dockerDesktopSocketExists = true
	}

	// Se o socket do Docker Desktop existe, não é modo nativo
	if dockerDesktopSocketExists {
		return false
	}

	// 5. Verifica se o Docker Engine está disponível
	// Tenta verificar se o socket do Docker Engine existe
	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		// Se o socket do Docker Engine não existe, não é modo nativo
		return false
	}

	// Se chegou aqui, está em Linux, não está em contêiner,
	// o socket do Docker Desktop não existe e o Docker Engine está disponível
	// Portanto, é modo nativo
	log.Log("- Detecção automática: Linux nativo detectado, ativando modo nativo")
	return true
}

// IsContainerMode detecta se está rodando em contêiner
func IsContainerMode() bool {
	return os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1"
}

// IsDockerDesktop detecta se está rodando com Docker Desktop
func IsDockerDesktop() bool {
	// Se estiver em modo contêiner ou modo nativo, não é Docker Desktop
	if IsContainerMode() || IsNativeMode() {
		return false
	}

	// Verifica se o socket do Docker Desktop existe
	if _, err := os.Stat("/run/host-services/backend.sock"); err == nil {
		return true
	}

	// Verifica paths alternativos do Docker Desktop em Linux
	paths := Paths()
	if paths.BackendSocket != "" {
		if _, err := os.Stat(paths.BackendSocket); err == nil {
			return true
		}
	}

	return false
}
