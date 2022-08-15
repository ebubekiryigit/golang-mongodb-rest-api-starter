package services

import (
	"errors"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNote create new note record
func CreateNote(userId primitive.ObjectID, title string, content string) (*db.Note, error) {
	note := db.NewNote(userId, title, content)
	err := mgm.Coll(note).Create(note)
	if err != nil {
		return nil, errors.New("cannot create new note")
	}

	return note, nil
}

// GetNotes get paginated note list
func GetNotes(userId primitive.ObjectID, page int, limit int) ([]db.Note, error) {
	var notes []db.Note

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Note{}).SimpleFind(
		&notes,
		bson.M{"author": userId},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find notes")
	}

	return notes, nil
}

func GetNoteById(userId primitive.ObjectID, noteId primitive.ObjectID) (*db.Note, error) {
	note := &db.Note{}
	err := mgm.Coll(note).First(bson.M{field.ID: noteId, "author": userId}, note)
	if err != nil {
		return nil, errors.New("cannot find note")
	}

	return note, nil
}

// UpdateNote updates a note with id
func UpdateNote(userId primitive.ObjectID, noteId primitive.ObjectID, request *models.NoteRequest) error {
	note := &db.Note{}
	err := mgm.Coll(note).FindByID(noteId, note)
	if err != nil {
		return errors.New("cannot find note")
	}

	if note.Author != userId {
		return errors.New("you cannot update this note")
	}

	note.Title = request.Title
	note.Content = request.Content
	err = mgm.Coll(note).Update(note)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteNote delete a note with id
func DeleteNote(userId primitive.ObjectID, noteId primitive.ObjectID) error {
	deleteResult, err := mgm.Coll(&db.Note{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: noteId, "author": userId})

	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete note")
	}

	return nil
}
