package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type nameBind struct {
	Name string `uri:"name"`
}

func (s *server) techs(c *gin.Context) {
	c.JSON(http.StatusOK, cs.StaticTechStore)
}

func (s *server) tech(c *gin.Context) {
	var name nameBind
	if err := c.ShouldBindUri(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tech := cs.StaticTechStore.GetTech(name.Name)
	if tech == nil {
		log.Error().Str("name", name.Name).Msg("get tech from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get tech from database"})
		return
	}

	c.JSON(http.StatusOK, tech)
}
