package service

import (
	"IM-Backend/service/mocks"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// encryptAES 使用 AES 的 CFB 模式对明文进行加密，并将结果进行 Base64 编码
func encryptAES(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	ciphertextWithIV := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(ciphertextWithIV), nil
}

func TestAuthSvc_decryptAES(t *testing.T) {
	type args struct {
		ciphertext string
		key        []byte
	}
	tests := []struct {
		name    string
		a       *AuthSvc
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid decryption",
			a:    &AuthSvc{},
			args: args{
				ciphertext: func() string {
					enc, _ := encryptAES("testdata", []byte("examplekey123456"))
					t.Log(enc)
					return enc
				}(),
				key: []byte("examplekey123456"),
			},
			want:    []byte("testdata"),
			wantErr: false,
		},
		{
			name: "invalid ciphertext",
			a:    &AuthSvc{},
			args: args{
				ciphertext: "invalidciphertext",
				key:        []byte("examplekey123456"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.decryptAES(tt.args.ciphertext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthSvc.decryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthSvc.decryptAES() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "valid verification",
			svc:  "testsvc",
			appKey: func() string {
				enc, err := encryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("testsecret666666"))
				if err != nil {
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
			name: "expired timestamp",
			svc:  "testsvc",
			appKey: func() string {
				enc, _ := encryptAES(fmt.Sprintf("%d", time.Now().Add(-20*time.Second).Unix()), []byte("testsecret666666"))
				return enc
			}(),
			secret: "testsecret666666",
			setupMock: func() {
				SvcHandlerMock.EXPECT().GetSecretByName("testsvc").Return("testsecret666666")
			},
			want: false,
		},
		{
			name: "invalid secret",
			svc:  "testsvc",
			appKey: func() string {
				enc, _ := encryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("nonosecret666666"))
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
