package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirgwain/craig-stars/game"
)

func (s *server) rules(c *gin.Context) {
	c.JSON(http.StatusOK, game.NewRules())
}
