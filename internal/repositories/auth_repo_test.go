package repositories

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "testing"

    "github.com/yemyoaung/managing-vehicle-tracking-models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    connStr = "mongodb://yoma_fleet:YomaFleet!123@localhost:27017"
)

func getAuthRepo() (*mongo.Client, *MongoAuthRepository, error) {
    // we can also use mock database for testing
    // but for now we will use real database to make sure everything is working fine
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
    if err != nil {
        return nil, nil, err
    }

    repo, err := NewMongoAuthRepository(context.Background(), client.Database("users"))

    if err != nil {
        return nil, nil, err
    }

    return client, repo, nil
}

func getRandomUser() (*models.User, error) {
    user, err := models.NewUser().SetEmail(fmt.Sprintf("hello%d@gmail.com", rand.Int()))
    if err != nil {
        return nil, err
    }

    if _, err = user.SetPassword(fmt.Sprintf("password%d", rand.Int())); err != nil {
        return nil, err
    }

    if _, err = user.SetRole(models.AdminRole); err != nil {
        return nil, err
    }

    if err = user.Build(); err != nil {
        return nil, err
    }
    return user, nil
}

func TestMongoAuthRepository_CreateUser(t *testing.T) {
    client, repo, err := getAuthRepo()

    if err != nil {
        t.Fatal(err)
    }

    defer func(client *mongo.Client, ctx context.Context) {
        err := client.Disconnect(ctx)
        if err != nil {
            log.Println("Failed to disconnect from database")
        }
    }(client, context.Background())

    user, err := getRandomUser()
    if err != nil {
        t.Fatal(err)
    }

    err = repo.CreateUser(context.Background(), user)

    if err != nil {
        t.Fatal(err)
    }

    if user.ID.IsZero() {
        t.Fatal("ID should not be zero")
    }

    err = repo.CreateUser(context.Background(), user)

    if err == nil {
        t.Fatal("Error should not be nil")
    }
}

func TestMongoAuthRepository_FindUserByEmail(t *testing.T) {

    client, repo, err := getAuthRepo()

    if err != nil {
        t.Fatal(err)
    }

    defer func(client *mongo.Client, ctx context.Context) {
        err := client.Disconnect(ctx)
        if err != nil {
            log.Println("Failed to disconnect from database")
        }
    }(client, context.Background())

    user, err := getRandomUser()

    if err != nil {
        t.Fatal(err)
    }

    err = repo.CreateUser(context.Background(), user)

    if err != nil {
        t.Fatal(err)
    }

    var dbUser1 models.User

    err = repo.FindAdminByEmail(context.Background(), string(user.Email), &dbUser1)

    if err != nil {
        t.Fatal(err)
        return
    }

    if err = dbUser1.Validate(); err != nil {
        t.Fatal(err)
    }

    if user.Email != dbUser1.Email {
        t.Fatal("Email should be equal")
    }

    var dbUser2 models.User

    err = repo.FindAdminByEmail(context.Background(), "radnomEmail@gmail.com", &dbUser2)

    if err == nil {
        t.Fatal("Error should not be nil")
    }

    if dbUser2.Check() == nil {
        t.Fatal("User should not be valid")
    }
}
