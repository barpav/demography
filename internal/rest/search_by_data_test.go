package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/demography/internal/rest/mocks"
	"github.com/barpav/demography/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_searchByData(t *testing.T) {
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
		wantBody    *models.SearchResultV1
		wantStatus  int
	}{
		{
			name: "Success with filters (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people?name=Ivan&limit=10", nil)
					r.Header.Set("Accept", models.MimeTypeSearchResultV1)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("SearchResultV1", mock.Anything, &models.SearchFilters{
						Name:  "Ivan",
						Limit: 10,
					}).Return(&models.SearchResultV1{
						Total: 2,
						Data: []*models.EnrichedPersonDataV1{
							{
								Id:         5,
								Surname:    "Ivanov",
								Name:       "Ivan",
								Patronymic: "Ivanovich",
							},
							{
								Id:      10,
								Surname: "Petrov",
								Name:    "Ivan",
							},
						},
					}, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": models.MimeTypeSearchResultV1,
			},
			wantBody: &models.SearchResultV1{
				Total: 2,
				Data: []*models.EnrichedPersonDataV1{
					{
						Id:         5,
						Surname:    "Ivanov",
						Name:       "Ivan",
						Patronymic: "Ivanovich",
					},
					{
						Id:      10,
						Surname: "Petrov",
						Name:    "Ivan",
					},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Success without filters (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people", nil)
					r.Header.Set("Accept", models.MimeTypeSearchResultV1)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("SearchResultV1", mock.Anything, mock.Anything).Return(
						&models.SearchResultV1{
							Total: 2,
							Data: []*models.EnrichedPersonDataV1{
								{
									Id:         5,
									Surname:    "Ivanov",
									Name:       "Ivan",
									Patronymic: "Ivanovich",
								},
								{
									Id:      10,
									Surname: "Petrov",
									Name:    "Ivan",
								},
							},
						}, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": models.MimeTypeSearchResultV1,
			},
			wantBody: &models.SearchResultV1{
				Total: 2,
				Data: []*models.EnrichedPersonDataV1{
					{
						Id:         5,
						Surname:    "Ivanov",
						Name:       "Ivan",
						Patronymic: "Ivanovich",
					},
					{
						Id:      10,
						Surname: "Petrov",
						Name:    "Ivan",
					},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Incorrect parameters (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people?age=old&limit=1000", nil)
					r.Header.Set("Accept", models.MimeTypeSearchResultV1)
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Requested media type is not supported (406)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/v1/people", nil)
					r.Header.Set("Accept", "application/xml")
					return r
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
			s.searchByData(tt.args.w, tt.args.r)

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

			var body *models.SearchResultV1
			decoded := models.SearchResultV1{}
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
