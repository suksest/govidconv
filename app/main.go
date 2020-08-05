package main

import (
	"os"

	"github.com/labstack/gommon/log"
	"gitlab.com/prosa-backend-test/video_converter/infra"
	"gitlab.com/prosa-backend-test/video_converter/video/handler"
	"gitlab.com/prosa-backend-test/video_converter/video/usecase"
)

func main() {
	err := os.Setenv("PORT", "5000")
	if err != nil {
		log.Error(err)
	}
	router := infra.NewServer()
	vu := usecase.NewVideoUsecase()
	handler.NewVideoHandler(router, vu)

	router.Static("static", "/app/files")
	router.Start(":" + os.Getenv("PORT"))
}
