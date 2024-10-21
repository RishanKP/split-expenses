package repository

import (
	"context"
	"split-expenses/library/utils"
	"split-expenses/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseRepository interface {
	CreateExpense(expense models.Expense) error
	GetAll(userId string) ([]models.Expense, error)
}

type repo struct {
	collection *mongo.Collection
}

func (r repo) CreateExpense(expense models.Expense) error {
	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.TODO(), expense)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) GetAll(userId string) ([]models.Expense, error) {
	filter := bson.M{}

	if !utils.IsEmpty(userId) {
		filter["participants.user_id"] = userId
	}

	curr, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return []models.Expense{}, err
	}

	var expenses []models.Expense

	for curr.Next(context.TODO()) {
		var e models.Expense
		err := curr.Decode(&e)
		if err != nil {
			return []models.Expense{}, err
		}

		expenses = append(expenses, e)
	}

	return expenses, nil
}

func NewExpenseRepository(db *mongo.Database, collectionName string) ExpenseRepository {
	return repo{
		collection: db.Collection(collectionName),
	}
}
