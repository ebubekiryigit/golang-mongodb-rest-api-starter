package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	mgm.DefaultModel `bson:",inline"`
	Author           primitive.ObjectID `json:"author" bson:"author"`
	Title            string             `json:"title" bson:"title"`
	Content          string             `json:"content" bson:"content"`
}

func NewNote(author primitive.ObjectID, title string, content string) *Note {
	return &Note{
		Author:  author,
		Title:   title,
		Content: content,
	}
}

func (model *Note) CollectionName() string {
	return "notes"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
