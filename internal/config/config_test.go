// Пакет config используется для чтения данных из файлов конфигурации
// и переменных окружения.

package config

import (
	"testing"
)

// TestMustLoad позволяет проверить корректность указания пути
// к файлу конфига в переменных окружения.

func TestMustLoad(t *testing.T) {
	var got *Config = MustLoad()
	if got == nil {
		t.Fatalf("MustLoad() error = failed to load config")
	}
	if got.Env != "local" && got.Env != "dev" && got.Env != "prod" {
		t.Errorf("MustLoad() error = failed to load config")
	}
}
