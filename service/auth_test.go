package service

import (
	"IM-Backend/pkg"
	"IM-Backend/service/mocks"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)


func TestAuthSvc_Verify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	SvcHandlerMock := mocks.NewMockSvcHandler(ctrl)

	tests := []struct {
		name      string
		svc       string
		appKey    string
		secret    string
		setupMock func()
		want      bool
	}{
		{
			name:   "valid verification",
			svc:    "testsvc",
			appKey: func() string {
				enc, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("testsecret666666"))
                if err!=nil {
                    t.Error(err)
                }
                t.Log(enc)
				return enc
			}(),
			secret: "testsecret666666",
			setupMock: func() {
				SvcHandlerMock.EXPECT().GetSecretByName("testsvc").Return("testsecret666666")
			},
			want: true,
		},
		{
			name:   "invalid appKey",
			svc:    "testsvc",
			appKey: "invalidappkey",
			secret: "testsecret666666",
			setupMock: func() {
				SvcHandlerMock.EXPECT().GetSecretByName("testsvc").Return("testsecret666666")
			},
			want: false,
		},
		{
			name:   "expired timestamp",
			svc:    "testsvc",
			appKey: func() string {
				enc, _ := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Add(-20*time.Second).Unix()), []byte("testsecret666666"))
				return enc
			}(),
			secret: "testsecret666666",
			setupMock: func() {
				SvcHandlerMock.EXPECT().GetSecretByName("testsvc").Return("testsecret666666")
			},
			want: false,
		},
		{
			name:   "invalid secret",
			svc:    "testsvc",
			appKey: func() string {
				enc, _ := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("nonosecret666666"))
				return enc
			}(),
			secret: "testsecret666666",
			setupMock: func() {
				SvcHandlerMock.EXPECT().GetSecretByName("testsvc").Return("testsecret666666")
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			a := &AuthSvc{svcHandler: SvcHandlerMock}
			if got := a.Verify(tt.svc, tt.appKey); got != tt.want {
				t.Errorf("AuthSvc.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}