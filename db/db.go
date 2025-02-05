package db

import (
	"context"
	"data-sender/config"
	"data-sender/core"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type dbconfig struct {
	timeZone              string
	sslMode               string
	poolMaxConns          int32
	poolMinConns          int32
	poolMaxConnLifetime   time.Duration
	poolMaxConnIdleTime   time.Duration
	poolHealthCheckPeriod time.Duration
}

type configOption func(cn *dbconfig)

func MinPoolConns(minConns int32) func(cn *dbconfig) {
	return func(c *dbconfig) {
		c.poolMinConns = minConns
	}
}

func MaxPoolConns(maxConns int32) func(cn *dbconfig) {
	return func(c *dbconfig) {
		c.poolMaxConns = maxConns
	}
}

func newDbConfig() dbconfig {
	return dbconfig{
		sslMode:               "disable",
		timeZone:              "UTC",
		poolMaxConns:          4,
		poolMinConns:          0,
		poolMaxConnLifetime:   time.Hour,
		poolMaxConnIdleTime:   time.Minute * 30,
		poolHealthCheckPeriod: time.Minute,
	}
}

func formatOption(url, option string, value interface{}) string {
	return url + " " + option + "=" + fmt.Sprintf("%v", value)
}

func addOptionsToConnStr(connStr string, options ...configOption) string {
	config := newDbConfig()
	for _, option := range options {
		option(&config)
	}

	connStr = formatOption(connStr, "sslmode", config.sslMode)
	connStr = formatOption(connStr, "TimeZone", config.timeZone)
	connStr = formatOption(connStr, "pool_max_conns", config.poolMaxConns)
	connStr = formatOption(connStr, "pool_min_conns", config.poolMinConns)
	connStr = formatOption(connStr, "pool_max_conn_lifetime", config.poolMaxConnLifetime)
	connStr = formatOption(connStr, "pool_max_conn_idle_time", config.poolMaxConnIdleTime)
	connStr = formatOption(connStr, "pool_health_check_period", config.poolHealthCheckPeriod)

	return connStr
}

func ConnectDb(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	log.Info().Str("host", cfg.Db.Host.Value).Str("name", cfg.Db.Name.Value).Msg("connecting to the database...")
	var err error

	if cfg.Db.Migrate.Value {
		log.Info().Msg("executing migrations")
		
		time.Sleep(5 * time.Second)
		if err = RunMigrations(cfg); err != nil {
			log.Warn().Err(err).Msg("error executing migrations")
			panic("migrations couldn't be ran")
		}
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Db.Host.Value, cfg.Db.Port.Value, cfg.Db.User.Value, cfg.Db.Pass.Value, cfg.Db.Name.Value)

	var pool *pgxpool.Pool

	url := addOptionsToConnStr(connStr, MinPoolConns(int32(cfg.Db.Pool.MinSize.Value)), MaxPoolConns(int32(cfg.Db.Pool.MaxSize.Value)))
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg(url)
	level, err := pgx.LogLevelFromString(cfg.Db.LogLevel.Value)
	if err != nil {
		return nil, err
	}
	poolConfig.ConnConfig.Logger = logger{level: level}

	for {
		pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err != nil {
			log.Error().Err(err).Msg("failed to create connection pool... retrying")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return pool, nil
}

type logger struct {
	level pgx.LogLevel
}

func (l logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if l.level < level {
		return
	}
	var evt *zerolog.Event
	switch level {
	case pgx.LogLevelTrace:
		evt = log.Trace()
	case pgx.LogLevelDebug:
		evt = log.Debug()
	case pgx.LogLevelInfo:
		evt = log.Info()
	case pgx.LogLevelWarn:
		evt = log.Warn()
	case pgx.LogLevelError:
		evt = log.Error()
	case pgx.LogLevelNone:
		evt = log.Info()
	default:
		evt = log.Info()
	}

	for k, v := range data {
		evt.Interface(k, v)
	}

	evt.Msg(msg)
}

func RunMigrations(cfg *config.Config) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User.Value,
		cfg.Db.Pass.Value,
		cfg.Db.Host.Value,
		cfg.Db.Port.Value,
		cfg.Db.Name.Value,
	)

	m, err := migrate.New("file:"+cfg.Db.MigrationFolder.Value, connStr)
	//m, err := migrate.New("file:///./db/migrations", connStr)
	log.Debug().Msg("запускаем миграции в папке " +  cfg.Db.MigrationFolder.Value)
	if err != nil {
		log.Error().Msg("в папке " +  cfg.Db.MigrationFolder.Value + " не были найдены файлы") 
		return err
	}
	if cfg.Db.Clean.Value {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Info().Msg("schema is up to date")
	}

	return nil
}

func GetQueryOptions(cn core.Conn, options ...core.QueryOptions) (conn core.Conn, forUpdate string) {
	conn = cn
	forUpdate = ""
	if len(options) > 0 {
		conn = options[0].Tx

		if options[0].ForUpdate {
			forUpdate = "FOR UPDATE"
		}
	}

	return conn, forUpdate
}

func GetUpdateOptions(cn core.Conn, options ...core.UpdateOptions) (conn core.Conn) {
	conn = cn
	if len(options) > 0 {
		conn = options[0].Tx
	}

	return conn
}