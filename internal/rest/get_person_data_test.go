package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/demography/internal/rest/mocks"
	"github.com/barpav/demography/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_getPersonData(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		testService testService
		args        args
		wantHeaders map[string]string
		wantBody    *models.EnrichedPersonDataV1
		wantStatus  int
	}{
		{
			name: "OK (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EnrichedPersonDataV1", mock.Anything, int64(101)).Return(
						&models.EnrichedPersonDataV1{
							Id:         101,
							Surname:    "Ivanov",
							Name:       "Ivan",
							Patronymic: "Ivanovich",
							Age:        50,
							Gender:     "male",
							Country:    "RU",
						},
						nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": models.MimeTypeEnrichedPersonDataV1,
			},
			wantBody: &models.EnrichedPersonDataV1{
				Id:         101,
				Surname:    "Ivanov",
				Name:       "Ivan",
				Patronymic: "Ivanovich",
				Age:        50,
				Gender:     "male",
				Country:    "RU",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Person not found in DB (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EnrichedPersonDataV1", mock.Anything, int64(101)).Return(nil, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Person not found - bad id (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "some-text")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Requested media type is not supported (406)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people/{id}", nil)
					r.Header.Set("Accept", "application/xml")
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.getPersonData(tt.args.w, tt.args.r)

			for k, v := range tt.wantHeaders {
				require.Equal(t, v, func() string {
					h := tt.args.w.Result().Header
					if h == nil {
						return ""
					}
					v := h[k]
					if len(v) == 0 {
						return ""
					}
					return v[0]
				}())
			}

			require.Equal(t, tt.wantStatus, tt.args.w.Code)

			if tt.wantBody == nil {
				return
			}

			var body *models.EnrichedPersonDataV1
			decoded := models.EnrichedPersonDataV1{}
			err := json.NewDecoder(tt.args.w.Body).Decode(&decoded)

			if err != nil && err != io.EOF {
				t.Fatal(err)
			}

			if err == nil {
				body = &decoded
			}

			require.Equal(t, body, tt.wantBody)
		})
	}
}
