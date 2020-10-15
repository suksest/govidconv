package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gitlab.com/prosa-backend-test/video_converter/domain"
	"gitlab.com/prosa-backend-test/video_converter/domain/mocks"
)

func TestNewVideoHandler(t *testing.T) {
	type args struct {
		e  *echo.Echo
		vu domain.VideoUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "new video handler",
			args: args{
				e:  echo.New(),
				vu: new(mocks.VideoUsecase),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewVideoHandler(tt.args.e, tt.args.vu)
		})
	}
}

func TestConvert(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/video/convert", strings.NewReader("test"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &VideoHandler{}

	// Assertions
	// test
	if assert.NoError(t, h.Convert(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"request Content-Type isn't multipart/form-data\"\n", rec.Body.String())
	}
}
