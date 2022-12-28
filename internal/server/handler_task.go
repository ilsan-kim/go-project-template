package server

import (
	"bytes"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"sampleProject/internal/params"
	"strconv"
	"time"
)

func (s *Server) taskPostHandler(c echo.Context) (err error) {
	ctx := context.Background()
	taskParam := new(params.CreateParams)

	if err = c.Bind(taskParam); err != nil {
		log.Println(err)
		return c.JSON(400, err.Error())
	}

	task, err := s.service.Task.Create(ctx, *taskParam)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	r := &Resp{
		Message: task,
	}
	return c.JSON(200, r)
}

func (s *Server) taskGetHandler(c echo.Context) (err error) {
	ctx := context.Background()
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	task, err := s.service.Task.Get(ctx, idInt)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	return c.JSON(200, task)
}

func (s *Server) taskGetListHandler(c echo.Context) (err error) {
	ctx := context.Background()
	tasks, err := s.service.Task.GetList(ctx)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	return c.JSON(200, tasks)
}

func (s *Server) taskGetCsvHandler(c echo.Context) (err error) {
	ctx := context.Background()
	data, _ := s.service.Task.MakeTaskListCsv(ctx)

	reader := bytes.NewReader(data)

	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; sample.csv")
	http.ServeContent(c.Response(), c.Request(), "sample.csv", time.Now(), reader)
	return nil
}

func (s *Server) taskDeleteHandler(c echo.Context) (err error) {
	ctx := context.Background()
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	err = s.service.Task.Delete(ctx, idInt)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, "ok")
}
