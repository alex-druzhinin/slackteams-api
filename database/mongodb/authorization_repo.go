package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"bitbucket.org/iwlab-standuply/slackteams-api/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	log "github.com/sirupsen/logrus"
)

const (
	authsCollectionName = "slack-bot-authorizations"
)

type slackBotAuthorizationsRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewSlackBotAuthorizationsRepository(uri string) handler.AuthorizationsRepository {
	clientOptions := options.Client()
	clientOptions.SetConnectTimeout(time.Duration(60) * time.Second)
	clientOptions.ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.WithError(err).Fatal()
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.WithError(err).Fatal()
	}

	connstr, err := connstring.Parse(uri)

	if err != nil {
		log.WithError(err).Fatal()
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.WithError(err).Error("NewSlackBotAuthorizationsRepository client.Ping")
	}

	db := client.Database(connstr.Database)

	return &slackBotAuthorizationsRepository{
		client,
		db,
	}
}

func (r *slackBotAuthorizationsRepository) GetAllAuthorizations(ctx context.Context) ([]*handler.SlackBotAuthorization, error) {
	filter := bson.D{{"enabled", true}}

	docs, err := r.findMany(ctx, filter)
	if err != nil {
		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	res := make([]*handler.SlackBotAuthorization, len(docs))

	for i, doc := range docs {
		res[i] = &handler.SlackBotAuthorization{
			AccessToken: doc.AccessToken,
			Scope:       doc.Scope,
			UserId:      doc.UserId,
			TeamName:    doc.TeamName,
			TeamId:      doc.TeamId,
			CreatedAt:   doc.CreatedAt.Format(time.RFC3339),
			Enabled:     doc.Enabled,
			Bot: handler.BotInfo{
				BotUserId:      doc.Bot.BotUserId,
				BotAccessToken: doc.Bot.BotAccessToken,
			},
		}
	}

	return res, nil
}

func (r *slackBotAuthorizationsRepository) GetAuthorization(ctx context.Context, teamId string) (*handler.SlackBotAuthorization, error) {
	filter := bson.D{{"enabled", true}, {"teamId", teamId}}

	doc, err := r.findOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, handler.ErrNotFound
		}
		
		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	res := &handler.SlackBotAuthorization{
		AccessToken: doc.AccessToken,
		Scope:       doc.Scope,
		UserId:      doc.UserId,
		TeamName:    doc.TeamName,
		TeamId:      doc.TeamId,
		CreatedAt:   doc.CreatedAt.Format(time.RFC3339),
		Enabled:     doc.Enabled,
		Bot: handler.BotInfo{
			BotUserId:      doc.Bot.BotUserId,
			BotAccessToken: doc.Bot.BotAccessToken,
		},
	}
	

	return res, nil
}

func (r *slackBotAuthorizationsRepository) findMany(ctx context.Context, filter interface{}) ([]*slackBotAuthorization, error) {
	var docs []*slackBotAuthorization

	collection := r.db.Collection(authsCollectionName)

	findOptions := options.Find()

	projection := bson.D{
		{"scope", 0},
	}

	findOptions.SetProjection(projection)

	cur, err := collection.Find(ctx, filter, findOptions)
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.WithError(err).Error("Failed to close cursor in findMany")
		}
	}()

	if err != nil {
		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	err = cur.All(ctx, &docs)
	if err != nil {
		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	log.WithContext(ctx).Debugf("findMany: %d\n", len(docs))
	log.WithContext(ctx).Debugf("findMany: %+v\n", docs[0])

	return docs, nil
}

func (r *slackBotAuthorizationsRepository) findOne(ctx context.Context, filter interface{}) (*slackBotAuthorization, error) {
	var doc *slackBotAuthorization

	collection := r.db.Collection(authsCollectionName)

	findOptions := options.FindOne()

	projection := bson.D{
		{"scope", 0},
	}

	findOptions.SetProjection(projection)

	err := collection.FindOne(ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	log.WithContext(ctx).Debugf("findOne: %+v\n", doc)

	return doc, nil
}
