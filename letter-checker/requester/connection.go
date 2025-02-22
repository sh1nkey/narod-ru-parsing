package requester

import (
	"letter-checker/config"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gocql/gocql"
)

type CheckerRepo interface {
	Init(cfg *config.Config)
	CheckExisted(letters string) bool
	Close()
}

type MockChecker struct {
	client map[string]bool
}

func (mc MockChecker) Init(cfg *config.Config) error {
	mc.client = make(map[string]bool)
	log.Info().Msgf("Создали моковое соединение к БД с url %s", cfg.Db.Host.Value)
	return nil
}

func (mc *MockChecker) CheckExisted(letters string) bool {
    if mc.client[letters] {
        return false
    }

    mc.client[letters] = true
   
    return true 
}

func (mc *MockChecker)Close() {}


type CassandraChecker struct {
	client *gocql.Session
}

func (cc *CassandraChecker) Init(cfg *config.Config) {
	host := cfg.Db.Host.Value
	log.Info().Msgf("Подключаемся к Cassandra с URL=%s", host)

	session := cc.createConn(host, cfg.Db.User.Value, cfg.Db.Pass.Value)
	cc.createNamespace(session)
	cc.createTable(session)

	cc.client = session
	
}

func (cc *CassandraChecker) CheckExisted(letters string) bool {
    query := "INSERT INTO cass_keyspace.cass_table (letters) VALUES (?) IF NOT EXISTS"

    var existingLetters string

    // Используем ScanCAS для обработки результата INSERT ... IF NOT EXISTS
    if applied, err := cc.client.Query(query, letters).ScanCAS(&existingLetters); err != nil {
        log.Err(err).Msgf("Ошибка при добавлении в БД letters=%s", letters)
        return false
    } else if !applied {
        log.Warn().Msgf("Такой уже есть в БД... letters=%s", existingLetters)
        return false
    }

    log.Info().Msgf("Буквы %s добавлены в БД :)", letters)
    return true
}


func (cc *CassandraChecker)Close() {
	cc.client.Close()
	log.Info().Msg("Клиент кассандры успешно закрыт")
}

func (cc *CassandraChecker) createConn(host string, username string, password string ) *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password, AllowedAuthenticators: []string {"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"}} //you will need to allow the use of the Instaclustr Password Authenticator.
	session, err := cluster.CreateSession()
        if err != nil {
		log.Fatal().Err(err).Msg("Couldn't connect to Cassandra")
	}
	return session
}

func (cc *CassandraChecker) createNamespace(sess *gocql.Session) {
	err := sess.Query("CREATE KEYSPACE IF NOT EXISTS cass_keyspace WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor': 1};").Exec() 
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось создать неймспейс")
	}
}

func (cc *CassandraChecker) createTable(sess *gocql.Session) {
	err := sess.Query("CREATE TABLE IF NOT EXISTS cass_keyspace.cass_table (letters text, PRIMARY KEY (letters));").Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось создать таблицу")
	}
}
