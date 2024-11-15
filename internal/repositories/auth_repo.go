package repositories

import (
    "context"
    "time"

    "github.com/yemyoaung/managing-vehicle-tracking-models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepository interface {
    CreateUser(ctx context.Context, user *models.User) error
    FindAdminByEmail(ctx context.Context, email string, user *models.User) error
    FindAdminByID(ctx context.Context, id string, user *models.User) error
}

type MongoAuthRepository struct {
    collection *mongo.Collection
}

func NewMongoAuthRepository(ctx context.Context, db *mongo.Database) (*MongoAuthRepository, error) {
    userCollection := db.Collection("users")

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    indexModel := mongo.IndexModel{
        Keys:    bson.M{"email": 1},
        Options: options.Index().SetUnique(true),
    }

    _, err := userCollection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return nil, err
    }
    return &MongoAuthRepository{
        collection: userCollection,
    }, nil
}

func (repo *MongoAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
    if err := user.Build(); err != nil {
        return err
    }
    result, err := repo.collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (repo *MongoAuthRepository) findOne(ctx context.Context, filter bson.M, user *models.User) error {
    err := repo.collection.FindOne(ctx, filter).Decode(user)
    if err != nil {
        return err
    }
    if err := user.Check(); err != nil {
        return err
    }
    return nil
}

func (repo *MongoAuthRepository) FindAdminByEmail(
    ctx context.Context,
    email string,
    user *models.User,
) error {
    return repo.findOne(ctx, bson.M{"email": email, "role": "admin"}, user)
}

func (repo *MongoAuthRepository) FindAdminByID(
    ctx context.Context,
    id string,
    user *models.User,
) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    return repo.findOne(ctx, bson.M{"_id": objID, "role": "admin"}, user)
}
