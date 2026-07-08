package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := `# comentário
PORT=9999
export STREAM_URL="https://stream.example/live"
NODE_API_KEY='chave secreta'
EMPTY_LINE_BELOW=ok

INVALIDO
DOTENV_PRESET=do-arquivo
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOTENV_PRESET", "do-ambiente")
	t.Setenv("PORT", "")
	t.Setenv("STREAM_URL", "")
	t.Setenv("NODE_API_KEY", "")

	loaded, err := LoadDotEnv(path)
	if err != nil || !loaded {
		t.Fatalf("loaded=%v err=%v", loaded, err)
	}
	if got := os.Getenv("PORT"); got != "9999" {
		t.Errorf("PORT=%q", got)
	}
	if got := os.Getenv("STREAM_URL"); got != "https://stream.example/live" {
		t.Errorf("STREAM_URL=%q (aspas duplas)", got)
	}
	if got := os.Getenv("NODE_API_KEY"); got != "chave secreta" {
		t.Errorf("NODE_API_KEY=%q (aspas simples)", got)
	}
	// variável já definida no ambiente NÃO é sobrescrita
	if got := os.Getenv("DOTENV_PRESET"); got != "do-ambiente" {
		t.Errorf("DOTENV_PRESET=%q, .env não deveria sobrescrever", got)
	}
}

func TestLoadDotEnvMissingFile(t *testing.T) {
	loaded, err := LoadDotEnv(filepath.Join(t.TempDir(), "nao-existe.env"))
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if loaded {
		t.Error("arquivo inexistente deveria retornar loaded=false")
	}
}
