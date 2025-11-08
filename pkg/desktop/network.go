package desktop

import (
	"os"
	"strings"

	"github.com/docker/mcp-gateway/pkg/log"
)

// GetHostAddress determina o endereço de host correto baseado no ambiente
// Usa a detecção automática do sistema e respeita override manual
func GetHostAddress() string {
	// 1. Respeita override manual via variável de ambiente
	if env := os.Getenv("DOCKER_MCP_NATIVE_MODE"); env != "" {
		if env == "1" {
			log.Log("- Usando modo nativo (override manual via DOCKER_MCP_NATIVE_MODE=1)")
			return "localhost"
		} else {
			log.Log("- Usando modo Docker Desktop (override manual via DOCKER_MCP_NATIVE_MODE!=1)")
			return "host.docker.internal"
		}
	}

	// 2. Usa detecção automática
	if IsNativeMode() {
		log.Log("- Detecção automática: Linux nativo detectado, usando localhost")
		return "localhost"
	}

	// 3. Para Docker Desktop ou outros ambientes
	log.Log("- Detecção automática: Docker Desktop ou outro ambiente detectado, usando host.docker.internal")
	return "host.docker.internal"
}

// GetHostAddressWithDefault retorna o endereço de host com um valor padrão opcional
// Se o endereço detectado for vazio, retorna o padrão especificado
func GetHostAddressWithDefault(defaultAddress string) string {
	address := GetHostAddress()
	if address == "" {
		log.Log("- Endereço detectado está vazio, usando padrão: " + defaultAddress)
		return defaultAddress
	}
	return address
}

// NormalizeHostAddress normaliza o endereço de host para garantir formato consistente
// Remove espaços em branco e converte para minúsculas
func NormalizeHostAddress(address string) string {
	return strings.TrimSpace(strings.ToLower(address))
}

// IsValidHostAddress verifica se o endereço de host é válido para os contextos suportados
func IsValidHostAddress(address string) bool {
	normalized := NormalizeHostAddress(address)
	validAddresses := []string{
		"localhost",
		"host.docker.internal",
		"127.0.0.1",
		"0.0.0.0",
	}

	for _, valid := range validAddresses {
		if normalized == valid {
			return true
		}
	}

	return false
}
