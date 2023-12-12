package adapters

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-management/app/courses/domain/model"
	"user-management/config"
)

type CategoryModel struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type ParentCourseModel struct {
	CourseId    primitive.ObjectID `json:"_id" bson:"_id"`
	CategoryId  string             `json:"category_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Tags        []string           `json:"tags"`
	IsPublish   int                `json:"is_publish"`
	SubCourses  []*CourseModel     `json:"sub_courses"`
}

type CourseModel struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Contents    []*ContentModel `json:"contents"`
}

type ContentModel struct {
	IsAssignment bool   `json:"is_assignment"`
	Ordering     int    `json:"ordering"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Content      string `json:"content"`
}

type FilterModel struct {
	ID         string `json:"id"`
	Filter     string `json:"filter"`
	CategoryId int    `json:"category_id"`
	Page       int64  `json:"page"`
	PerPage    int64  `json:"per_page"`
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

func (c *CourseMongoRepository) courseCollection() *mongo.Collection {
	return c.mongoDatabase.Collection(c.cfg.CourseCollection)
}

func (c *CourseMongoRepository) categoryCollection() *mongo.Collection {
	return c.mongoDatabase.Collection(c.cfg.CategoryCollection)
}

func (c *CourseMongoRepository) AddCategory(ctx context.Context, categoryName string) error {
	coll := c.categoryCollection()
	_, err := coll.InsertOne(ctx, bson.D{{"$inc", bson.D{{"category_id", 1}}}, {"category_name", categoryName}})
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *CourseMongoRepository) GetCategories(ctx context.Context) ([]*CategoryModel, error) {
	out := make([]*CategoryModel, 0)
	coll := c.categoryCollection()
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		c.log.Error(err)
		return out, err
	}

	for cursor.Next(ctx) {
		data := &CategoryModel{}
		err = cursor.Decode(data)
		if err != nil {
			c.log.Error(err)
			continue
		}
		out = append(out, data)
	}
	return out, nil
}

func (c *CourseMongoRepository) AddCourse(ctx context.Context, cr *model.CoursePath) error {
	collection := c.courseCollection()
	dataCr, err := marshalCourse(cr)
	if err != nil {
		c.log.Error(err)
		return err
	}

	dataCr.CourseId = primitive.NewObjectID()
	_, err = collection.InsertOne(ctx, cr)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *CourseMongoRepository) GetCourses(ctx context.Context, in *FilterModel) ([]*model.CoursePath, error) {
	out := make([]*model.CoursePath, 0)
	query := make([]bson.M, 0)
	findQuery := make([]bson.M, 0)
	if in.ID != "" {
		query = append(query, bson.M{"_id": in.ID})
	}
	if in.CategoryId != 0 {
		query = append(query, bson.M{"category_id": in.CategoryId})
	}
	if in.Filter != "" {
		query = append(query, bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: in.Filter, Options: ""}}})
		query = append(query, bson.M{"tags": bson.M{"$regex": primitive.Regex{Pattern: in.Filter, Options: ""}}})
	}
	findQuery = append(findQuery, bson.M{"$match": func() bson.M {
		if len(query) > 0 {
			return bson.M{"$and": query}
		}
		return bson.M{}
	}})
	script, _ := bson.Marshal(findQuery)
	c.log.Debugf("filter bson: %s", string(script))

	limit := in.PerPage
	if limit == 0 {
		limit = 10
	}

	page := in.Page
	if page == 0 {
		page = 1
	}

	skip := limit * (page - 1)
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}
	coll := c.courseCollection()
	cursor, err := coll.Find(ctx, findQuery, &fOpt)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	for cursor.Next(ctx) {
		each := &ParentCourseModel{}
		if err = cursor.Decode(each); err != nil {
			c.log.Error(err)
			continue
		}
		eachCr, err := unmarshalCourse(each)
		if err != nil {
			c.log.Error(err)
			continue
		}
		if !eachCr.IsCourseVisible() {
			continue
		}
		out = append(out, eachCr)
	}

	return out, nil
}

func (c *CourseMongoRepository) UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, cm *model.CoursePath) (*model.CoursePath, error)) error {
	coll := c.courseCollection()
	existing := &ParentCourseModel{}
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(existing)
	if err != nil {
		c.log.Error(err)
		return err
	}
	existingPath, err := unmarshalCourse(existing)
	if err != nil {
		c.log.Error(err)
		return err
	}

	updateCourse, err := updateFn(ctx, existingPath)
	if err != nil {
		c.log.Error(err)
		return err
	}
	_, err = coll.ReplaceOne(ctx, bson.D{{"_id", id}}, updateCourse)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *CourseMongoRepository) ReplaceCourse(ctx context.Context, id string, in *model.CoursePath) error {
	coll := c.courseCollection()
	existing := &ParentCourseModel{}
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(existing)
	if err != nil {
		c.log.Error(err)
		return err
	}
	_, err = coll.ReplaceOne(ctx, bson.D{{"_id", id}}, in)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func getNextSequence(collectionName string) {

}

func marshalCourse(in *model.CoursePath) (*ParentCourseModel, error) {
	out := &ParentCourseModel{}
	err := mapstructure.Decode(in, out)
	if err != nil {
		return nil, err
	}
	if in.CourseId != "" {
		out.CourseId, _ = primitive.ObjectIDFromHex(in.CourseId)
	}
	return out, nil
}

func unmarshalCourse(in *ParentCourseModel) (*model.CoursePath, error) {
	out := &model.CoursePath{}
	err := mapstructure.Decode(in, out)
	if err != nil {
		return nil, err
	}

	out.CourseId = in.CourseId.String()
	return out, nil
}
