package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/demography/internal/rest/mocks"
	"github.com/barpav/demography/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_editPersonData(t *testing.T) {
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
		wantStatus  int
	}{
		{
			name: "Edited (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", models.MimeTypeEditedPersonDataV1)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UpdatePersonDataV1", mock.Anything, int64(101),
						&models.EditedPersonDataV1{
							Name:       "Ivan",
							Patronymic: "Ivanovich",
							Surname:    "Ivanov",
						},
					).Return(nil)
					return s
				}(),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Incomplete person data (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", models.MimeTypeEditedPersonDataV1)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Incorrect person data (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := struct{ Test string }{Test: "test value"}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", models.MimeTypeEditedPersonDataV1)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Person not found in DB (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", models.MimeTypeEditedPersonDataV1)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("UpdatePersonDataV1", mock.Anything, int64(101),
						&models.EditedPersonDataV1{
							Name:       "Ivan",
							Patronymic: "Ivanovich",
							Surname:    "Ivanov",
						},
					).Return(ErrPersonDataNotFoundTest{})
					return s
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Person not found - bad id (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", models.MimeTypeEditedPersonDataV1)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "some-text")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Unsupported person data (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.EditedPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("PUT", "/v1/people/{id}", &buf)
					r.Header.Set("Content-Type", "application/json")
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "some-text")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantStatus: http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.editPersonData(tt.args.w, tt.args.r)

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}

type ErrPersonDataNotFoundTest struct{}

func (e ErrPersonDataNotFoundTest) Error() string {
	return "person data not found (test)"
}

func (e ErrPersonDataNotFoundTest) ImplementsPersonDataNotFoundError() {
}
