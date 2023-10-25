package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/demography/internal/rest/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_deletePersonData(t *testing.T) {
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
			name: "Deleted (204)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("DeletePersonData", mock.Anything, int64(101)).Return(nil)
					return s
				}(),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Person not found in DB (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "101")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("DeletePersonData", mock.Anything, int64(101)).Return(ErrPersonDataNotFoundTest{})
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
					r := httptest.NewRequest("DELETE", "/v1/people/{id}", nil)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "some-text")
					return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.deletePersonData(tt.args.w, tt.args.r)

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}
