package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type UsersHandler struct {
	UserUsecase domain.UsersUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(e *echo.Echo, us domain.UsersUsecase) {
	handler := &UsersHandler{
		UserUsecase: us,
	}
	e.GET("/users", handler.Fetch)
	//e.POST("/users", handler.Store)
	e.GET("/users/:id", handler.GetByID)
	e.DELETE("/users/delete/:id", handler.Delete)
}

// FetchArticle will fetch the article based on given params
func (a *UsersHandler) Fetch(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		size, _ := strconv.Atoi(c.QueryParam("size"))
		ctx := c.Request().Context()

		art, err := a.UserUsecase.Fetch(ctx,page,size)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, art)
}

// GetByID will get article by given id
func (a *UsersHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	art, err := a.UserUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the article by given request body
//func (a *ArticleHandler) Store(c echo.Context) (err error) {
//	var article domain.Article
//	err = c.Bind(&article)
//	if err != nil {
//		return c.JSON(http.StatusUnprocessableEntity, err.Error())
//	}
//
//	var ok bool
//	if ok, err = isRequestValid(&article); !ok {
//		return c.JSON(http.StatusBadRequest, err.Error())
//	}
//
//	ctx := c.Request().Context()
//	err = a.AUsecase.Store(ctx, &article)
//	if err != nil {
//		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
//	}
//
//	return c.JSON(http.StatusCreated, article)
//}
//
//// Delete will delete article by given param
func (a *UsersHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	err := a.UserUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
