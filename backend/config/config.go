// 浣滅敤锛氶厤缃粨鏋勪綋涓庡姞杞介€昏緫锛岀敤浜庝粠閰嶇疆鏂囦欢鍜岀幆澧冨彉閲忎腑璇诲彇搴旂敤閰嶇疆
package config

import (
	"log"
	"strings"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-clean-architecture/internal/pkg/logger"
)

var ProviderSet = wire.NewSet(NewConfig, ProvideLoggerOptions)

type Config struct {
	Env      string         `mapstructure:"env"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type PostgresConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NewConfig 浣跨敤 viper 浠庢枃浠跺拰鐜鍙橀噺鍔犺浇閰嶇疆
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config.default")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // 鏀寔浠庢墽琛岀洰褰曚笅鐨?config 璇诲彇
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	
	// 鍏佽鐜鍙橀噺瑕嗙洊
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 灏濊瘯璇诲彇 config.default.yaml
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Viper read config warning: %v", err)
	}

	// 鍙€夛細灏濊瘯鍚堝苟 config.local.yaml
	v.SetConfigName("config.local")
	if err := v.MergeInConfig(); err != nil {
		// local configuration is optional, don't throw an error if it doesn't exist
		log.Printf("Local config not merged: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 璁剧疆榛樿鐨勬棩蹇楀€硷紙闃叉閰嶇疆娌″～瀵艰嚧鏃犳硶杈撳嚭锛?
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.Filename == "" {
		cfg.Logger.Filename = "logs/app.log"
	}
	if cfg.Logger.MaxSize == 0 {
		cfg.Logger.MaxSize = 100 // 榛樿 100MB
	}
	if cfg.Logger.MaxBackups == 0 {
		cfg.Logger.MaxBackups = 3
	}
	if cfg.Logger.MaxAge == 0 {
		cfg.Logger.MaxAge = 7
	}

	return &cfg, nil
}

// ProvideLoggerOptions 灏?Config 杞寲涓?Logger 鎵€闇€鐨?Options
func ProvideLoggerOptions(cfg *Config) *logger.Options {
	return &logger.Options{
		Level:      cfg.Logger.Level,
		Filename:   cfg.Logger.Filename,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
	}
}
