// Пакет config используется для чтения данных из файлов конфигурации
// и переменных окружения.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Структура конфиг файла
type Config struct {
	Env         string `yaml:"env" env-default:"prod"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	//
	// Закомментированные поля ниже понадобятся позже, при развертывании
	// на удаленной машине.
	// StorageUser   string `yaml:"storage_user" env:"DB_USER" env-required:"true"`
	// StoragePasswd string `yaml:"storage_passwd" env:"DB_PASSWD" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}
type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"0.0.0.0:80"`
	AddressTLS   string        `yaml:"addressTLS" env-default:"0.0.0.0:443"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad - инициализирует данные из конфиг файла. Путь к файлу берет из
// переменной окружения RS_CONFIG_PATH. Если не удается, то завершает
// приложение с ошибкой.
func MustLoad() *Config {
	configPath := os.Getenv("RS_CONFIG_PATH")

	if configPath == "" {
		log.Fatal("RS_CONFIG_PATH is not set")
	}

	// Проверяем, существует ли файл конфига
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
