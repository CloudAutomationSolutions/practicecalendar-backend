package model

import (
	"context"

	"cloud.google.com/go/firestore"
)

// Datastore ...
// Main interface to implement for the store where we keep the User and Project data
type Datastore interface {
	SetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

// DB ...
// Main data structure holding the object interacting with the database.
// Will implement the Datastore interface
type DB struct {
	Client *firestore.Client
}

// NewDB ...
// DB Object Builder
func NewDB(ctx context.Context, dataSourceName string) (*DB, error) {
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		return nil, err
	}

	return &DB{client}, nil
}

// GetUser ...
// Get user from Datastore
func (d *DB) GetUser(ctx context.Context, id string) (*User, error) {

	userDoc := d.Client.Doc("Users/" + id)
	docsnap, err := userDoc.Get(ctx)
	if err != nil {
		return nil, err
	}

	var userData User
	if err := docsnap.DataTo(&userData); err != nil {
		return nil, err
	}

	return &userData, nil
}

// SetUser ...
// Create or update existing User in database
func (d *DB) SetUser(ctx context.Context, user *User) error {

	userDoc := d.Client.Doc("Users/" + user.ID)
	_, err := userDoc.Set(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
