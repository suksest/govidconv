package handler

import (
	"testing"

	"github.com/labstack/echo"
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
