package repository

import (
	"context"
	"time"

	"split-expenses/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetById(ctx context.Context, id string) (models.User, error)
}

type userRepo struct {
	collection *mongo.Collection
}

func (r userRepo) Create(ctx context.Context, user models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)

	return err
}

func (r userRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {

	filter := bson.M{"email": email}

	var user models.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r userRepo) GetById(ctx context.Context, id string) (models.User, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	filter := bson.M{"_id": objId}

	var user models.User
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func NewUserRepository(db *mongo.Database, collectionName string) UserRepository {
	collection := db.Collection(collectionName)
	return userRepo{
		collection: collection,
	}
}
