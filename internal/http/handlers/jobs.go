package handlers

import (
	"net/http"

	"github.com/icloudeng/platform-installer/internal/resources/db"

	"github.com/gin-gonic/gin"
)

type (
	jobIdUri struct {
		ID uint `uri:"id" binding:"required,number"`
	}

	jobsHandler struct{}
)

var Jobs jobsHandler

func (jobsHandler) GetJobsByID(c *gin.Context) {
	var data jobIdUri

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	job := db.Jobs.JobGetByID(data.ID)

	if job == nil {
		c.AbortWithStatusJSON(404, gin.H{"msg": "Job Not Found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
