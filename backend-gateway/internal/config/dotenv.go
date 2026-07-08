package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv carrega pares KEY=VALUE de um arquivo .env para o ambiente do
// processo, sem sobrescrever variáveis já definidas. Suporta comentários (#),
// linhas em branco, prefixo "export " e valores entre aspas simples/duplas.
// Arquivo inexistente não é erro (retorna false).
func LoadDotEnv(path string) (loaded bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return false, err
		}
	}
	return true, scanner.Err()
}
