package controllers

import (
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

// CreateNewNote godoc
// @Summary      Create Note
// @Description  creates a new note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        req  body      models.NoteRequest true "Note Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /notes [post]
// @Security     ApiKeyAuth
func CreateNewNote(c *gin.Context) {
	var requestBody models.NoteRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	note, err := services.CreateNote(userId.(primitive.ObjectID), requestBody.Title, requestBody.Content)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"note": note}
	response.SendResponse(c)
}

// GetNotes godoc
// @Summary      Get Notes
// @Description  gets user notes with pagination
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        page  query    string  false  "Switch page by 'page'"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /notes [get]
// @Security     ApiKeyAuth
func GetNotes(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	notes, _ := services.GetNotes(userId.(primitive.ObjectID), page, limit)
	hasPrev := page > 0
	hasNext := len(notes) > limit

	if hasNext {
		notes = notes[:limit]
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"notes": notes, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

// GetOneNote godoc
// @Summary      Get a note
// @Description  get note by id
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /notes/{id} [get]
// @Security     ApiKeyAuth
func GetOneNote(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	noteId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	note, err := services.GetNoteFromCache(userId.(primitive.ObjectID), noteId)
	if err == nil {
		models.SendResponseData(c, gin.H{"note": note, "cache": true})
		return
	}

	note, err = services.GetNoteById(userId.(primitive.ObjectID), noteId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	go services.CacheOneNote(userId.(primitive.ObjectID), note)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"note": note}
	response.SendResponse(c)
}

// UpdateNote godoc
// @Summary      Update a note
// @Description  updates a note by id
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id     path    string  true  "Note ID"
// @Param        req    body    models.NoteRequest true "Note Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /notes/{id} [put]
// @Security     ApiKeyAuth
func UpdateNote(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	noteId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	var noteRequest models.NoteRequest
	_ = c.ShouldBindBodyWith(&noteRequest, binding.JSON)

	err := services.UpdateNote(userId.(primitive.ObjectID), noteId, &noteRequest)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

// DeleteNote godoc
// @Summary      Delete a note
// @Description  deletes note by id
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /notes/{id} [delete]
// @Security     ApiKeyAuth
func DeleteNote(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	noteId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	err := services.DeleteNote(userId.(primitive.ObjectID), noteId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}
