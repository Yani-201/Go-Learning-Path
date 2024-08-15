package tests

import (
    "context"
    "os"
    "testing"
    "time"

    "task-manager-api-clean/domain"
    "task-manager-api-clean/repository"

    "github.com/stretchr/testify/suite"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    nonExistentId  = "nonExistentId"
)

type TaskRepositoryTestSuite struct {
    suite.Suite
    repo       domain.TaskRepository
    database   *mongo.Database
    collection string
    client     *mongo.Client
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    suite.NoError(err)

    suite.client = client
    suite.database = client.Database("test_db")
    suite.collection = "tasks"
    suite.repo = repository.NewTaskRepository(suite.database, suite.collection)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
    if suite.client != nil {
        err := suite.client.Disconnect(context.Background())
        suite.NoError(err)
    }
}

func (suite *TaskRepositoryTestSuite) SetupTest() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    suite.database.Collection(suite.collection).Drop(ctx)
}

func (suite *TaskRepositoryTestSuite) TestCreate_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Title: "Test Task", Description: "Description of test task", Status: "Pending", DueDate: time.Now()}
    createdTask, err := suite.repo.Create(ctx, task)
    suite.NoError(err)
    suite.Equal(task.Title, createdTask.Title)

    var result domain.Task
    err = suite.database.Collection(suite.collection).FindOne(ctx, bson.M{"_id": createdTask.Id}).Decode(&result)
    suite.NoError(err)
    suite.Equal(task.Title, result.Title)
}

func (suite *TaskRepositoryTestSuite) TestUpdate_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Title: "Test Task", Description: "Description of test task", Status: "Pending", DueDate: time.Now()}
    suite.repo.Create(ctx, task)

    updateTask := &domain.Task{Title: "Updated Task Title", Description: "Updated Description"}
    updatedTask, err := suite.repo.Update(ctx, task.Id, updateTask)
    suite.NoError(err)
    suite.Equal(updateTask.Title, updatedTask.Title)

    var result domain.Task
    err = suite.database.Collection(suite.collection).FindOne(ctx, bson.M{"_id": task.Id}).Decode(&result)
    suite.NoError(err)
    suite.Equal(updateTask.Title, result.Title)
}

func (suite *TaskRepositoryTestSuite) TestUpdate_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Attempting to update a non-existent task
    _, err := suite.repo.Update(ctx, nonExistentId, &domain.Task{Title: "New Title"})
    suite.Error(err)
}

func (suite *TaskRepositoryTestSuite) TestUpdate_NoFieldsUpdated() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Title: "Test Task", Description: "Description of test task", Status: "Pending", DueDate: time.Now()}
    suite.repo.Create(ctx, task)

    updateTask := &domain.Task{Title: "Test Task", Description: "Description of test task"}
    _, err := suite.repo.Update(ctx, task.Id, updateTask)
    suite.Error(err)
    suite.Equal("task not updated, no new information is provided", err.Error())
}

func (suite *TaskRepositoryTestSuite) TestDelete_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Title: "Test Task", Description: "Description of test task", Status: "Pending", DueDate: time.Now()}
    suite.repo.Create(ctx, task)

    err := suite.repo.Delete(ctx, task.Id)
    suite.NoError(err)

    var result domain.Task
    err = suite.database.Collection(suite.collection).FindOne(ctx, bson.M{"_id": task.Id}).Decode(&result)
    suite.Error(err)
}

func (suite *TaskRepositoryTestSuite) TestDelete_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Attempting to delete a non-existent task
    err := suite.repo.Delete(ctx, nonExistentId)
    suite.Error(err)
}

func (suite *TaskRepositoryTestSuite) TestGetAll_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task1 := &domain.Task{Title: "Test Task 1", Description: "Description of test task 1", Status: "Pending", DueDate: time.Now()}
    task2 := &domain.Task{Title: "Test Task 2", Description: "Description of test task 2", Status: "Completed", DueDate: time.Now()}
    suite.repo.Create(ctx, task1)
    suite.repo.Create(ctx, task2)

    tasks, err := suite.repo.GetAll(ctx)
    suite.NoError(err)
    suite.Len(*tasks, 2)
}

func (suite *TaskRepositoryTestSuite) TestGetAll_NoTasks() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tasks, err := suite.repo.GetAll(ctx)
    suite.NoError(err)
    suite.Len(*tasks, 0)
}

func (suite *TaskRepositoryTestSuite) TestGetById_Success() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    task := &domain.Task{Title: "Test Task", Description: "Description of test task", Status: "Pending", DueDate: time.Now()}
    suite.repo.Create(ctx, task)

    fetchedTask, err := suite.repo.GetById(ctx, task.Id)
    suite.NoError(err)
    suite.Equal(task.Title, fetchedTask.Title)
}

func (suite *TaskRepositoryTestSuite) TestGetById_Failure() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Attempting to fetch a non-existent task
    _, err := suite.repo.GetById(ctx, nonExistentId)
    suite.Error(err)
}

func TestTaskRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(TaskRepositoryTestSuite))
}








