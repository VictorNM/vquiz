package main

import (
	"github.com/labstack/echo"
	_ "github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	db *mongo.Database
	e  *echo.Echo
}

func (s *server) routes() {
	s.e.GET("/questions", s.handleViewQuestions())
	s.e.GET("/questions/add", s.handleAddQuestionGET())
	s.e.POST("/questions/add", s.handleAddQuestionPOST())
}

func (s *server) handleViewQuestions() echo.HandlerFunc {
	return func(c echo.Context) error {
		questions, err := getQuestions(c.Request().Context(), s.db)
		if err != nil {
			return c.Render(200, "view-questions", nil)
		}

		return c.Render(200, "view-questions", questions)
	}
}

func (s *server) handleAddQuestionGET() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, "add-question", nil)
	}
}

func (s *server) handleAddQuestionPOST() echo.HandlerFunc {
	return func(c echo.Context) error {
		q := Question{
			Content: c.FormValue("question"),
			Answer:  c.FormValue("question"),
		}

		if err := insertQuestion(c.Request().Context(), s.db, q); err != nil {
			return c.Render(500, "add-question", nil)
		}


		return c.Render(200, "add-question", nil)
	}
}
