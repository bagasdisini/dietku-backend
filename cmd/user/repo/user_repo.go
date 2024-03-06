package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	FullName  string             `json:"fullname" bson:"fullname"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
}

type Users []User

func DecodeAsUsers(cursor *mongo.Cursor) (*Users, error) {
	docs := Users{}
	err := cursor.All(context.TODO(), &docs)
	if err != nil {
		return nil, err
	}
	return &docs, nil
}

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		coll: db.Collection("users"),
	}
}

func (r *UserRepository) FindOne(id primitive.ObjectID) (*User, error) {
	var d = &User{}
	err := r.coll.FindOne(context.TODO(), bson.M{"_id": id, "isDeleted": bson.M{"$ne": true}}).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *UserRepository) FindOneByEmail(email string) (*User, error) {
	var d = &User{}
	err := r.coll.FindOne(context.TODO(), bson.M{"email": email, "isDeleted": bson.M{"$ne": true}}).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *UserRepository) InsertOne(newUser *User) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(context.TODO(), newUser)
}

func (r *UserRepository) UpdateOne(user *User) (*User, error) {
	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": user,
	}

	var d = &User{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.coll.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
