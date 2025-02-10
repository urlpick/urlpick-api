package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	urls := r.Group("/urls")
	{
		urls.POST("", h.CreateURL)
		urls.GET("/:hash", h.GetURL)
	}
}

func (h *Handler) CreateURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("invalid request body"))
		return
	}

	resp, err := h.service.CreateShortURL(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetURL(c *gin.Context) {
	hash := c.Param("hash")

	resp, err := h.service.GetURL(c.Request.Context(), hash)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
