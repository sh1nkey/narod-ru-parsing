package requester

import (
	"context"
	"letter-checker/config"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CheckerRepo interface {
	Init(cfg *config.Config)
	ReadOne(letters string) bool
	WriteOne(letters string)
}









type MongoChecker struct {
	client *mongo.Client
	collection *mongo.Collection
}

func (mnc *MongoChecker) Init(cfg *config.Config) {
	mnc.createConnection(cfg.Db.Host.Value, cfg.Db.User.Value, cfg.Db.Pass.Value)
	mnc.createTable()
}

func (mnc *MongoChecker) ReadOne(letters string) bool {
	doc := LettersDTO{Letters: letters}

	log.Info().Msgf("делаем запрос в монгу для буков %s", letters)

	filter := bson.M{"letters": letters}
	log.Info().Msgf("получили фильтр для монги %s", filter)

	err := mnc.collection.FindOne(context.TODO(), filter).Decode(&doc)

	if err != nil {
		log.Err(err).Msg("Не нашли документ")
		return false
	}
	
	log.Info().Msg("Нашли документ")
	return true
}

func (mnc *MongoChecker) WriteOne(letters string) {

	doc := LettersDTO{Letters: letters}
	insertResult, err := mnc.collection.InsertOne(context.TODO(), &doc)

	if err != nil {
		log.Err(err).Msgf("Не смогли добавить в монгу запись %s", letters)
		return
	}
	
	log.Info().Msgf("Успешно добавили документ: %s", insertResult.InsertedID)
}




func (mnc *MongoChecker) createConnection(host string, username string, password string) {
    // Строка подключения с указанием authSource
	log.Info().Msgf("%s %s", username, password)
    url := "mongodb://" + username + ":" + password + "@" + host + ":27017/?authSource=admin"

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
    if err != nil {
        log.Fatal().Err(err).Msg("Не можем подключится к mongo")
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal().Err(err).Msg("Не смогли пингануть монгу")
    }

    log.Info().Msg("Connected to MongoDB!")
    mnc.client = client
}


func (mnc *MongoChecker) createTable() {
	collection := mnc.client.Database("mongo").Collection("letters")
	mnc.collection = collection

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "letters", Value: 1}},
        Options: options.Index().SetUnique(true), 
    }
  
    if _, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
        log.Fatal().Err(err).Msg("Не создать индекс")
    }
}
