package repository_test

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "task-manager-api-clean/domain"
    "task-manager-api-clean/repository"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    repo       domain.UserRepository
    database   *mongo.Database
    collection string
    client     *mongo.Client
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    suite.NoError(err)

    suite.client = client
    suite.database = client.Database("test_db")
    suite.collection = "users"
    suite.repo = repository.NewUserRepository(suite.database, suite.collection)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
    if suite.client != nil {
        err := suite.client.Disconnect(context.Background())
        suite.NoError(err)
    }
}

func (suite *UserRepositoryTestSuite) SetupTest() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    suite.database.Collection(suite.collection).Drop(ctx)
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    user := &domain.User{Username: "test", Password: "test", Email: "test@example.com"}
    createdUser, err := suite.repo.Create(ctx, user)
    suite.NoError(err)
    suite.Equal(user.Username, createdUser.Username)

    var result domain.User
    err = suite.database.Collection(suite.collection).FindOne(ctx, bson.M{"username": "test"}).Decode(&result)
    suite.NoError(err)
    suite.Equal(user.Username, result.Username)
	suite.Equal(user.Email, result.Email)
}


func (suite *UserRepositoryTestSuite) TestGetByUsername_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    user := &domain.User{Username: "test3", Password: "test3", Email: "test3@example.com"}
    suite.repo.Create(ctx, user)

    fetchedUser, err := suite.repo.GetByUsername(ctx, "test3")
    suite.NoError(err)
    suite.Equal(user.Username, fetchedUser.Username)
}

func (suite *UserRepositoryTestSuite) TestGetByUsername_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := suite.repo.GetByUsername(ctx, "nonExistentUser")
    suite.Error(err)
}

func (suite *UserRepositoryTestSuite) TestUpdateRole_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	user1 := &domain.User{Username: "test0", Password: "test0", Email: "test0@example.com"}
    suite.repo.Create(ctx, user1)

    user := &domain.User{Username: "test4", Password: "test4", Email: "test4@example.com"}
    suite.repo.Create(ctx, user)

    updatedUser, err := suite.repo.UpdateRole(ctx, "test4", "admin")
    suite.NoError(err)
    suite.Equal("admin", updatedUser.Role)
}

func (suite *UserRepositoryTestSuite) TestUpdateRole_negative() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	user1 := &domain.User{Username: "testnew", Password: "testnew", Email: "testnew@example.com"}
    suite.repo.Create(ctx, user1)

    user := &domain.User{Username: "test5", Password: "test5", Email: "test5@example.com"}
    suite.repo.Create(ctx, user)

    // First promotion to admin
    updatedUser, err := suite.repo.UpdateRole(ctx, "test5", "admin")
    suite.NoError(err)
    suite.Equal("admin", updatedUser.Role)

    // Attempt to promote again to admin
    updatedUser, err = suite.repo.UpdateRole(ctx, "test5", "admin")
    suite.Error(err)
    suite.Nil(updatedUser)
    suite.EqualError(err, "user is already an admin")
}

func (suite *UserRepositoryTestSuite) TestUpdateRole_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := suite.repo.UpdateRole(ctx, "nonExistentUser", "admin")
    suite.Error(err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}