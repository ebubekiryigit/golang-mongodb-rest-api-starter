package validators

import (
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"net/http"
)

func CreateNoteValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createNoteRequest models.NoteRequest
		_ = c.ShouldBindBodyWith(&createNoteRequest, binding.JSON)

		if err := createNoteRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func GetNotesValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.DefaultQuery("page", "0")
		err := validation.Validate(page, is.Int)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "invalid page: "+page)
			return
		}

		c.Next()
	}
}

func UpdateNoteValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var noteRequest models.NoteRequest
		_ = c.ShouldBindBodyWith(&noteRequest, binding.JSON)

		if err := noteRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
