package config

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)



const (
	AppName  = "NarodRu Parser"
	Revision = "1"
)

var (
	// Build time arguments
	AppVersion  string
	Sha1Version string
	BuildTime   string

	// Runtime flags
	profile      *string
	port         *string
	configSource *string
	configUrl    *string
	configBranch *string
	configUser   *string
	configPass   *string
)

type StringConfig struct {
	Value       string `json:"value"   yaml:"value"`
	Default     string `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
}

type BoolConfig struct {
	Value       bool   `json:"value"   yaml:"value"`
	Default     bool   `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
}

type IntConfig struct {
	Value       int64  `json:"value"   yaml:"value"`
	Default     int64  `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
}

type FloatConfig struct {
	Value       float64 `json:"value"   yaml:"value"`
	Default     float64 `json:"default" yaml:"default"`
	Description string  `json:"description" yaml:"description"`
}

type Config struct {
	AppName     StringConfig `json:"appName"     yaml:"appName"`
	AppVersion  StringConfig `json:"appVersion"  yaml:"appVersion"`
	Sha1Version StringConfig `json:"sha1Version" yaml:"sha1Version"`
	BuildTime   StringConfig `json:"buildTime"   yaml:"buildTime"`
	Profile     StringConfig `json:"profile"     yaml:"profile"`
	Revision    StringConfig `json:"revision"    yaml:"revision"`
	Port        StringConfig `json:"port"        yaml:"port"`
	Config      ConfigSource `json:"config"      yaml:"config"`
	Log         LogConfig    `json:"log"         yaml:"log"`
	Db          DbConfig     `json:"db"          yaml:"db"`
}

type ConfigSource struct {
	Print       BoolConfig   `json:"print"  yaml:"print"`
	Source      StringConfig `json:"source" yaml:"source"`
	Description string       `json:"description" yaml:"description"`
}


type LogConfig struct {
	Level       StringConfig `json:"level"      yaml:"level"`
	Structured  BoolConfig   `json:"structured" yaml:"structured"`
	Description string       `json:"description" yaml:"description"`
}

type DbConfig struct {
	Name            StringConfig `json:"name"            yaml:"name"`
	Host            StringConfig `json:"host"            yaml:"host"`
	Port            StringConfig `json:"port"            yaml:"port"`
	Migrate         BoolConfig   `json:"migrate"         yaml:"migrate"`
	MigrationFolder StringConfig `json:"migrationFolder" yaml:"migrationFolder"`
	Clean           BoolConfig   `json:"clean"           yaml:"clean"`
	User            StringConfig `json:"user"            yaml:"user"`
	Pass            StringConfig `json:"pass"            yaml:"pass"`
	Pool            DbPoolConfig `json:"pool"            yaml:"pool"`
	LogLevel        StringConfig `json:"logLevel"        yaml:"logLevel"`
	Description     string       `json:"description"     yaml:"description"`
}

type DbPoolConfig struct {
	MinSize           IntConfig `json:"minPoolSize"       yaml:"minPoolSize"`
	MaxSize           IntConfig `json:"maxPoolSize"       yaml:"maxPoolSize"`
	MaxConnLife       IntConfig `json:"maxConnLife"       yaml:"maxConnLife"`
	MaxConnIdle       IntConfig `json:"maxConnIdle"       yaml:"maxConnIdle"`
	HealthCheckPeriod IntConfig `json:"healthCheckPeriod" yaml:"healthCheckPeriod"`
	Description       string    `json:"description" yaml:"description"`
}

func (c *Config) Print() {
	if c.Config.Print.Value {
		log.Info().Interface("config", c).Msg("the following configurations have successfully loaded")
	}
}

func init() {
	def := &Config{}
	setupDefaults(def)

	profile = flag.String("p", def.Profile.Default, def.Profile.Description)
	port = flag.String("port", def.Port.Default, def.Port.Description)
	configSource = flag.String("s", def.Config.Source.Default, def.Config.Source.Description)

	viper.SetDefault("port", def.Port.Default)
	viper.SetDefault("profile", def.Profile.Default)

	viper.SetDefault("config.print", def.Config.Print.Default)
	viper.SetDefault("config.source", def.Config.Source.Default)

	viper.SetDefault("log.level", def.Log.Level.Default)
	viper.SetDefault("log.structured", def.Log.Structured.Default)

	viper.SetDefault("db.name", def.Db.Name.Default)
	viper.SetDefault("db.host", def.Db.Host.Default)
	viper.SetDefault("db.port", def.Db.Port.Default)
	viper.SetDefault("db.user", def.Db.User.Default)
	viper.SetDefault("db.pass", def.Db.Pass.Default)
	viper.SetDefault("db.clean", def.Db.Clean.Default)
	viper.SetDefault("db.migrate", def.Db.Migrate.Default)
	viper.SetDefault("db.migrationFile", def.Db.MigrationFolder.Default)
	viper.SetDefault("db.pool.minSize", def.Db.Pool.MinSize.Default)
	viper.SetDefault("db.pool.maxSize", def.Db.Pool.MaxSize.Default)

	
}


func setupDefaults(config *Config) {
	config.AppName = StringConfig{Value: AppName, Default: AppName, Description: "Name of the application in a human readable format. Example: Go Micro Example"}

	config.AppVersion = StringConfig{Value: AppVersion, Default: "", Description: "Semantic version of the application. Example: v1.2.3"}
	config.Sha1Version = StringConfig{Value: Sha1Version, Default: "", Description: "Git sha1 hash of the application version."}
	config.BuildTime = StringConfig{Value: BuildTime, Default: "", Description: "When this version of the application was compiled."}
	config.Profile = StringConfig{Value: "local", Default: "local", Description: "Running profile of the application, can assist with sensible defaults or change behavior. Examples: local, dev, prod"}
	config.Revision = StringConfig{Value: Revision, Default: Revision, Description: "A hard coded revision handy for quickly determining if local changes are running. Examples: 1, Two, 9999"}
	config.Port = StringConfig{Value: "8080", Default: "8080", Description: "Port that the application will bind to on startup. Examples: 8080, 3000"}

	config.Config.Description = "Settings for where and how the application should get its configurations."
	config.Config.Print = BoolConfig{Value: false, Default: false, Description: "Print configurations on startup."}
	config.Config.Source = StringConfig{Value: "", Default: "", Description: "Where the application should go for configurations. Examples: local, etcd"}

	config.Log.Description = "Settings for applicaton logging."
	config.Log.Level = StringConfig{Value: "trace", Default: "trace", Description: "The lowest level that the application should log at. Examples: info, warn, error."}
	config.Log.Structured = BoolConfig{Value: false, Default: false, Description: "Whether the application should output structured (json) logging, or human friendly plain text."}

	config.Db.Description = "Database configurations."
	config.Db.Name = StringConfig{Value: "micro-ex-db", Default: "micro-ex-db", Description: "The name of the database to connect to."}
	config.Db.Host = StringConfig{Value: "5432", Default: "5432", Description: "Port of the database."}
	config.Db.Migrate = BoolConfig{Value: true, Default: true, Description: "Whether or not database migrations should be executed on startup."}
	config.Db.MigrationFolder = StringConfig{Value: "db/migrations/", Default: "db/migrations", Description: "Location of migration files to be executed on startup."}
	config.Db.Clean = BoolConfig{Value: false, Default: false, Description: "WARNING: THIS WILL DELETE ALL DATA FROM THE DB. Used only during migration. If clean is true, all 'down' migrations are executed."}
	config.Db.User = StringConfig{Value: "postgres", Default: "postgres", Description: "User the application will use to connect to the database."}
	config.Db.Pass = StringConfig{Value: "postgres", Default: "postgres", Description: "Password the application will use for connecting to the database."}
	config.Db.Pool.MinSize = IntConfig{Value: 1, Default: 1, Description: "The minimum size of the pool."}
	config.Db.Pool.MaxSize = IntConfig{Value: 3, Default: 3, Description: "The maximum size of the pool."}
	config.Db.Pool.MaxConnLife = IntConfig{Value: time.Hour.Milliseconds(), Default: time.Hour.Milliseconds(), Description: "The maximum time a connection can live in the pool in milliseconds."}
	config.Db.Pool.MaxConnIdle = IntConfig{Value: time.Minute.Milliseconds() * 30, Default: time.Minute.Milliseconds() * 30, Description: "The maximum time a connection can idle in the pool in milliseconds."}
	config.Db.LogLevel = StringConfig{Value: "trace", Default: "trace", Description: "The logging level for database interactions. See: log.level"}
}

func loadLocalConfigs(filename string, config *Config) error {
	log.Info().Msg("loading local configurations...")

	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../.")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		ValueToConfigValue(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
	if err != nil {
		return err
	}

	return nil
}


func Load(filename string) *Config {
	config := &Config{}
	setupDefaults(config)

	var err error
	switch config.Config.Source.Value {
	case "local":
		err = loadLocalConfigs(filename, config)
	case "etcd":
		err = loadRemoteConfigs(config)
	default:
		log.Warn().
			Str("configSource", config.Config.Source.Value).
			Msg("unrecognized configuration source, using local")

		err = loadLocalConfigs(filename, config)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}

	err = loadCommandLineOverrides(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}

	return config
}

func loadRemoteConfigs(config *Config) error {
	return nil
}


func loadCommandLineOverrides(config *Config) error {
	flag.Parse()
	if *profile != config.Profile.Default {
		log.Debug().Str("profile", *profile).Str("config.Profile", config.Profile.Value).Str("config.Profile.Default", config.Profile.Default).Msg("overriding profile")
		config.Profile.Value = *profile
	}
	if *port != config.Port.Default {
		log.Debug().Str("port", *port).Str("config.port", config.Port.Value).Msg("overriding port")
		config.Port.Value = *port
	}
	if *configSource != config.Config.Source.Default {
		log.Debug().Str("configSource", *configSource).Str("config.Config.Source", config.Config.Source.Value).Msg("overriding config source")
		config.Config.Source.Value = *configSource
	}
	return nil
}


func ValueToConfigValue() mapstructure.DecodeHookFunc {
	return func(f reflect.Value, t reflect.Value) (interface{}, error) {

		if t.Kind() != reflect.Struct {
			return f.Interface(), nil
		}

		to := t.Interface()
		switch t := to.(type) {
		case IntConfig:
			v, err := getInt(f)
			if err != nil {
				return nil, err
			}
			t.Value = v
			return t, nil
		case StringConfig:
			v, err := getString(f)
			if err != nil {
				return nil, err
			}
			t.Value = v
			return t, nil
		case BoolConfig:
			v, err := getBool(f)
			if err != nil {
				return nil, err
			}
			t.Value = v
			return t, nil
		case FloatConfig:
			v, err := getFloat(f)
			if err != nil {
				return nil, err
			}
			t.Value = v
			return t, nil
		}

		return f.Interface(), nil
	}
}

func getString(f reflect.Value) (string, error) {
	data := f.Interface()

	switch f.Kind() {
	case reflect.Int64:
		raw := data.(int64)
		return strconv.FormatInt(raw, 10), nil
	case reflect.Int:
		raw := data.(int)
		return strconv.Itoa(raw), nil
	case reflect.String:
		raw := data.(string)
		return raw, nil
	case reflect.Bool:
		raw := data.(bool)
		return strconv.FormatBool(raw), nil
	case reflect.Float64:
		raw := data.(float64)
		return strconv.FormatFloat(raw, 'f', 3, 64), nil
	}

	return "", errors.New("unrecognized type")
}

func getBool(f reflect.Value) (bool, error) {
	data := f.Interface()

	switch f.Kind() {
	case reflect.Int64:
		raw := data.(int64)
		return raw > 0, nil
	case reflect.Int:
		raw := data.(int)
		return raw > 0, nil
	case reflect.String:
		raw := data.(string)
		return raw == "true", nil
	case reflect.Bool:
		return data.(bool), nil
	case reflect.Float64:
		raw := data.(float64)
		return raw > 0, nil
	}

	return false, errors.New("unrecognized type")
}

func getFloat(f reflect.Value) (float64, error) {
	data := f.Interface()

	switch f.Kind() {
	case reflect.Int64:
		raw := data.(int64)
		return float64(raw), nil
	case reflect.Int:
		raw := data.(int)
		return float64(raw), nil
	case reflect.String:
		raw := data.(string)
		return strconv.ParseFloat(raw, 64)
	case reflect.Bool:
		raw := data.(bool)
		if raw {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.Float64:
		raw := data.(float64)
		return raw, nil
	}

	return -1, errors.New("unrecognized type")
}

func getInt(f reflect.Value) (int64, error) {
	data := f.Interface()

	switch f.Kind() {
	case reflect.Int64:
		raw := data.(int64)
		return raw, nil
	case reflect.Int:
		raw := data.(int)
		return int64(raw), nil
	case reflect.String:
		raw := data.(string)
		return strconv.ParseInt(raw, 10, 64)
	case reflect.Bool:
		raw := data.(bool)
		if raw {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.Float64:
		raw := data.(float64)
		return int64(raw), nil
	}

	return -1, errors.New("unrecognized type")
}
