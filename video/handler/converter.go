package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/prosa-backend-test/video_converter/domain"
)

//VideoHandler ...
type VideoHandler struct {
	VideoUsecase domain.VideoUsecase
}

//NewVideoHandler ...
func NewVideoHandler(e *echo.Echo, vu domain.VideoUsecase) {
	handler := &VideoHandler{
		VideoUsecase: vu,
	}
	e.POST("/video/convert", handler.Convert)
}

func parseOutputPreset(c echo.Context) (op domain.OutputPreset) {
	compression := c.FormValue("compression")
	format := c.FormValue("format")
	vidCodec := c.FormValue("video_codec")
	audCodec := c.FormValue("audio_codec")

	if compression == "" {
		compression = "1.0"
	}
	if format == "" {
		format = "keep"
	}
	if vidCodec == "" {
		vidCodec = "default"
	}
	if audCodec == "" {
		audCodec = "default"
	}

	op.Compression = compression
	op.Format = format
	op.VideoCodec = vidCodec
	op.AudioCodec = audCodec
	return
}

//Convert ...
func (vh *VideoHandler) Convert(c echo.Context) (err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	op := parseOutputPreset(c)

	downloadPath, err := vh.VideoUsecase.Convert(file, op)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"downloadPath": c.Request().Host + "/" + downloadPath,
	})
}
