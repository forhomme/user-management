package adapters

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"user-management/config"
)

type AuditTrailModel struct {
	UserId   string
	Menu     string
	Method   string
	Request  string
	Response string
}

type CourseMongoRepository struct {
	cfg           *config.Config
	log           logger.Logger
	mongoDatabase *mongo.Database
}

func NewCourseMongoRepository(cfg *config.Config, log logger.Logger, mongoDatabase *mongo.Database) *CourseMongoRepository {
	return &CourseMongoRepository{
		cfg:           cfg,
		log:           log,
		mongoDatabase: mongoDatabase,
	}
}

func (c *CourseMongoRepository) auditTrailCollection() *mongo.Collection {
	return c.mongoDatabase.Collection(c.cfg.AuditTrailCollection)
}

func (c *CourseMongoRepository) AddAuditTrail(ctx context.Context, in *AuditTrailModel) error {
	coll := c.auditTrailCollection()
	_, err := coll.InsertOne(ctx, in)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}
