package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sampleProject/internal/common"
	"sampleProject/internal/service/task"
	"time"
)

type Resp struct {
	Message interface{} `json:"message"`
}

type Service struct {
	Task *task.TaskService
}

type Server struct {
	service    *Service
	httpServer *http.Server

	IsRunning bool
}

func New(service *Service) *Server {
	server := new(Server)

	// conf 에서 서비스를 주입받지 않고 그냥 인자로 받아서 직접 초기화함
	server.service = service

	e := echo.New()
	e.Match([]string{"GET"}, "/hello", server.helloHandler)
	e.Match([]string{"POST"}, "/tasks", server.taskPostHandler)
	e.Match([]string{"GET"}, "/tasks", server.taskGetListHandler)
	e.Match([]string{"GET"}, "/tasks/:id", server.taskGetHandler)
	e.Match([]string{"DELETE"}, "/tasks/:id", server.taskDeleteHandler)
	e.Match([]string{"GET"}, "/tasks/export", server.taskGetCsvHandler)

	srv := &http.Server{
		Addr:    ":8001",
		Handler: e,
	}
	server.httpServer = srv
	return server
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutdownGraceFully(ctx context.Context) error {
	s.IsRunning = false
	err := s.httpServer.Shutdown(ctx)
	for {
		cnt := len(common.GracefulShutdownJob)
		if cnt == 0 {
			break
		}
		log.Println("remained job : ", cnt)
		f, _ := os.OpenFile("jobs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		pprof.Lookup("goroutine").WriteTo(f, 1)
		f.Close()
		time.Sleep(time.Second * 5)
	}
	return err
}
