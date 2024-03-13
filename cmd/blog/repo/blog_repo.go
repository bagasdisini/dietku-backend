package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type By struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	FullName string             `json:"fullname" bson:"fullname"`
	At       time.Time          `json:"at" bson:"at"`
}

type Blog struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Header    string             `json:"header" bson:"header"`
	Content   string             `json:"content" bson:"content"`
	Category  []string           `json:"category" bson:"category"`
	CreatedBy By                 `json:"createdBy" bson:"createdBy"`
	UpdatedBy *By                `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
}

type Blogs []Blog

func DecodeAsBlogs(cursor *mongo.Cursor) (*Blogs, error) {
	docs := Blogs{}
	err := cursor.All(context.TODO(), &docs)
	if err != nil {
		return nil, err
	}
	return &docs, nil
}

type BlogRepository struct {
	coll *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		coll: db.Collection("blogs"),
	}
}

func (r *BlogRepository) FindAll() (*Blogs, error) {
	cursor, err := r.coll.Find(context.TODO(), bson.M{"isDeleted": bson.M{"$ne": true}})
	if err != nil {
		return nil, err
	}
	return DecodeAsBlogs(cursor)
}

func (r *BlogRepository) FindOne(id primitive.ObjectID) (*Blog, error) {
	var d = &Blog{}
	err := r.coll.FindOne(context.TODO(), bson.M{"_id": id, "isDeleted": bson.M{"$ne": true}}).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *BlogRepository) FindByCategory(category string) (*Blogs, error) {
	cursor, err := r.coll.Find(context.TODO(), bson.M{"category": category, "isDeleted": bson.M{"$ne": true}})
	if err != nil {
		return nil, err
	}
	return DecodeAsBlogs(cursor)
}

func (r *BlogRepository) FindByUser(userID primitive.ObjectID) (*Blogs, error) {
	cursor, err := r.coll.Find(context.TODO(), bson.M{"createdBy._id": userID, "isDeleted": bson.M{"$ne": true}})
	if err != nil {
		return nil, err
	}
	return DecodeAsBlogs(cursor)
}

func (r *BlogRepository) InsertOne(newBlog *Blog) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(context.TODO(), newBlog)
}

func (r *BlogRepository) UpdateOne(blog *Blog) (*Blog, error) {
	filter := bson.M{"_id": blog.ID}

	update := bson.M{
		"$set": blog,
	}

	var d = &Blog{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.coll.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *BlogRepository) DeleteOne(id primitive.ObjectID) (*Blog, error) {
	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{"isDeleted": true},
	}

	var d = &Blog{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.coll.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
