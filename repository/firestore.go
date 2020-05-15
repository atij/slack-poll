package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/atij/slack-poll/model"
)

// FirestoreRepository ...
type FirestoreRepository struct {
	client *firestore.Client
	collection string
}

// NewFirestoreRepository ...
func NewFirestoreRepository(c *firestore.Client, col string) (*FirestoreRepository, error) {
	return &FirestoreRepository{
		client: c,
		collection: col,
	}, nil
}

// Create ...
func (repo *FirestoreRepository) Create(p *model.Poll) error {
	ref := repo.client.Collection(repo.collection).NewDoc()
	p.ID = ref.ID
	_, err := ref.Set(context.Background(), p)

	if err!= nil {
		return err
	}
	
	return nil
}

// Update ...
func (repo *FirestoreRepository) Update(id string, p *model.Poll) error {
	_, err := repo.client.Collection(repo.collection).Doc(id).Set(context.Background(), &p)
	if err != nil {
		return err
	}
	return nil
}

// Find ...
func (repo *FirestoreRepository) Find(id string) (*model.Poll, error) {
	dsnap, err := repo.client.Collection(repo.collection).Doc(id).Get(context.Background())
	
	if err != nil {
		return nil, err
	}

	var p model.Poll
	dsnap.DataTo(&p)

	return &p, nil
}
