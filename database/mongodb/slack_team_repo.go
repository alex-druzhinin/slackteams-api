package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"bitbucket.org/iwlab-standuply/slackteams-api/handler"
	"bitbucket.org/iwlab-standuply/slackteams-api/rpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	log "github.com/sirupsen/logrus"
)

const (
	slackTeamsCollectionName = "slack-teams"
)

type slackTeamsRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewSlackTeamsRepository(uri string) rpc.SlackTeamsRepository {
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
		log.WithError(err).Error("NewslackTeamsRepository client.Ping")
	}

	db := client.Database(connstr.Database)

	return &slackTeamsRepository{
		client,
		db,
	}
}

func (r *slackTeamsRepository) FindTeamByID(ctx context.Context, teamId string) (*rpc.SlackTeam, error) {
	filter := bson.D{{"id", teamId}}

	doc, err := r.findOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, handler.ErrNotFound
		}

		log.WithContext(ctx).WithError(err).Debug()
		return nil, err
	}

	res := &rpc.SlackTeam{
		ID: doc.TeamID,
		Name: doc.Name,
		Domain: doc.Domain,
		EmailDomain: doc.EmailDomain,
		Icon: rpc.SlackIcon{
			Image34: doc.Icon.Image34,
			Image44: doc.Icon.Image44,
			Image68: doc.Icon.Image68,
			Image88: doc.Icon.Image88,
			Image102: doc.Icon.Image102,
			Image132: doc.Icon.Image132,
			Image230: doc.Icon.Image230,
			ImageDefault: doc.Icon.ImageDefault,
		},
		IsDeleted: doc.IsDeleted,
		DeletedAt: doc.DeletedAt,
		CreatedAt: doc.CreatedAt,
		Tags: doc.Tags,
	}

	return res, nil
}

func (r *slackTeamsRepository) findOne(ctx context.Context, filter interface{}) (*slackTeam, error) {
	var doc *slackTeam

	collection := r.db.Collection(slackTeamsCollectionName)

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
