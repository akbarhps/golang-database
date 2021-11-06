package repository

import (
	"context"
	"golang-database/entity"
	"testing"

	golang_database "golang-database"
)

func TestCommentInsert(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment := &entity.Comment{
		Email:   "test@test.com",
		Comment: "test comment dengan repository",
	}

	comment, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindById(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment, err := commentRepository.FindById(ctx, 1)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindByIdNotFound(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment, err := commentRepository.FindById(ctx, 0)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindAll(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comments, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		t.Logf("%+v", comment)
	}
}
