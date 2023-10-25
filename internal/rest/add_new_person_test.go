package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/demography/internal/rest/mocks"
	"github.com/barpav/demography/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_addNewPerson(t *testing.T) {
	type testService struct {
		cfg     *config
		stats   StatisticsProvider
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
			name: "New person data added (201)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/v1/people", &buf)
					r.Header.Set("Content-Type", models.MimeTypeNewPersonDataV1)
					return r
				}(),
			},
			testService: testService{
				cfg: &config{statsTimeout: 3000},
				stats: func() *mocks.StatisticsProvider {
					s := mocks.NewStatisticsProvider(t)
					s.On("AgeByName", "Ivan").Return(50, nil)
					s.On("GenderByName", "Ivan").Return("male", nil)
					s.On("CountryByName", "Ivan").Return("RU", nil)
					return s
				}(),
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("CreateNewPersonDataV1", mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type": models.MimeTypeEnrichedPersonDataV1,
			},
			wantBody: &models.EnrichedPersonDataV1{
				Surname:    "Ivanov",
				Name:       "Ivan",
				Patronymic: "Ivanovich",
				Age:        50,
				Gender:     "male",
				Country:    "RU",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Incomplete new person data (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/v1/people", &buf)
					r.Header.Set("Content-Type", models.MimeTypeNewPersonDataV1)
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Unsupported new person data (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					m := models.NewPersonDataV1{
						Name:       "Ivan",
						Patronymic: "Ivanovich",
						Surname:    "Ivanov",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(m)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/v1/people", &buf)
					r.Header.Set("Content-Type", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:     tt.testService.cfg,
				stats:   tt.testService.stats,
				storage: tt.testService.storage,
			}
			s.addNewPerson(tt.args.w, tt.args.r)

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
