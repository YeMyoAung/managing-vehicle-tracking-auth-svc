package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/repositories"
    "github.com/yemyoaung/managing-vehicle-tracking-models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var args = []string{
    "MongoDB URL",
    "Email",
    "Password",
}

func main() {
    var dbUrl, email, password string
    var err error
    for v := range args {
        log.Println("What is the", args[v], "?")
        if v == 0 {
            _, err = fmt.Scan(&dbUrl)
            continue
        }
        if v == 1 {
            _, err = fmt.Scan(&email)
            continue
        }
        if v == 2 {
            _, err = fmt.Scan(&password)
            continue
        }
    }

    // Connect to MongoDB
    db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbUrl))
    if err != nil {
        log.Println("Failed to connect to MongoDB", err)
        return
    }

    repo, err := repositories.NewMongoAuthRepository(context.Background(), db.Database("users"))

    if err != nil {
        log.Println("Failed to create repository", err)
        return
    }

    user := models.NewUser()

    _, err = user.SetEmail(email)
    if err != nil {
        log.Println("Failed to set email", err)
        return
    }

    _, err = user.SetPassword(password)
    if err != nil {
        log.Println("Failed to set password", err)
        return
    }

    _, err = user.SetRole(models.AdminRole)

    if err != nil {
        log.Println("Failed to set role", err)
        return
    }

    err = repo.CreateUser(context.Background(), user)

    if err != nil {
        log.Println("Failed to create user", err)
        return
    }

    log.Println("User created successfully")
}
