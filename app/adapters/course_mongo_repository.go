package adapters

import (
	"context"
	"errors"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"go.opentelemetry.io/otel/codes"
	"user-management/app/domain/course"

	"github.com/forhomme/app-base/usecase/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-management/config"
)

type CourseMongoRepository struct {
	cfg           *config.Config
	log           logger.Logger
	tracer        *telemetry.OtelSdk
	mongoDatabase *mongo.Database
}

func NewCourseMongoRepository(cfg *config.Config, log logger.Logger, mongoDatabase *mongo.Database, tracer *telemetry.OtelSdk) *CourseMongoRepository {
	return &CourseMongoRepository{
		cfg:           cfg,
		log:           log,
		tracer:        tracer,
		mongoDatabase: mongoDatabase,
	}
}

func (c *CourseMongoRepository) courseCollection() *mongo.Collection {
	return c.mongoDatabase.Collection(c.cfg.CourseCollection)
}

func (c *CourseMongoRepository) categoryCollection() *mongo.Collection {
	return c.mongoDatabase.Collection(c.cfg.CategoryCollection)
}

func (c *CourseMongoRepository) AddCategory(ctx context.Context, categoryName string) error {
	return errors.New("not implemented")
}

func (c *CourseMongoRepository) GetCategories(ctx context.Context) ([]*course.Category, error) {
	return nil, errors.New("not implemented")
}

func (c *CourseMongoRepository) AddCourse(ctx context.Context, cr *course.CoursePath) (err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.add_course")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	collection := c.courseCollection()
	cr.CourseId = primitive.NewObjectID().String()
	_, err = collection.InsertOne(ctx, cr)
	if err != nil {
		return err
	}
	return nil
}

func (c *CourseMongoRepository) GetCourses(ctx context.Context, in *course.FilterCourse) (out []*course.CoursePath, err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.get_course")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	out = make([]*course.CoursePath, 0)
	query := make([]bson.M, 0)
	findQuery := bson.M{}
	if in.ID != "" {
		objectId, _ := primitive.ObjectIDFromHex(in.ID)
		query = append(query, bson.M{"_id": objectId})
	}
	if in.CategoryId != 0 {
		query = append(query, bson.M{"category.category_id": in.CategoryId})
	}
	if in.Filter != "" {
		query = append(query, bson.M{"title": bson.D{{"$regex", primitive.Regex{Pattern: in.Filter, Options: ""}}}})
		query = append(query, bson.M{"tags": bson.D{{"$regex", primitive.Regex{Pattern: in.Filter, Options: ""}}}})
	}

	if len(query) > 0 {
		findQuery = bson.M{"$or": query}
	}

	limit := in.PerPage
	page := in.Page

	skip := limit * (page - 1)
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}
	coll := c.courseCollection()
	cursor, err := coll.Find(ctx, findQuery, &fOpt)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		each := &ParentCourseModel{}
		if err = cursor.Decode(each); err != nil {
			c.log.Error(err)
			continue
		}
		eachCp := &course.CoursePath{
			IsPublished:    each.IsPublished,
			TotalSubCourse: each.TotalSubCourse,
			Category:       each.Category,
			CourseId:       each.CourseId.Hex(),
			Title:          each.Title,
			Description:    each.Description,
			SubCourses:     each.SubCourses,
		}
		out = append(out, eachCp)
	}

	return out, nil
}

func (c *CourseMongoRepository) UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context,
	cm *course.CoursePath) (*course.CoursePath, error)) (err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.update_course")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	coll := c.courseCollection()
	existing := &course.CoursePath{}
	err = coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(existing)
	if err != nil {
		return err
	}

	updateCourse, err := updateFn(ctx, existing)
	if err != nil {
		return err
	}
	_, err = coll.ReplaceOne(ctx, bson.D{{"_id", id}}, updateCourse)
	if err != nil {
		return err
	}
	return nil
}
