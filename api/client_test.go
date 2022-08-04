package api

import (
	mockdb "bank-mvp/db/mock"
	db "bank-mvp/db/sqlc"
	"bank-mvp/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateClientAPI(t *testing.T) {
	client := randomClient()
	testCases := []struct {
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"first_name": client.FirstName,
				"last_name": client.LastName,
				"country_id": client.CountryID,
				"user_id": client.UserID,
			},
			buildStubs: func (store *mockdb.MockStore) {
				arg := db.CreateClientParams {
					FirstName: client.FirstName,
					LastName: client.LastName,
					CountryID: client.CountryID,
					UserID: client.UserID,
				}

				store.EXPECT().CreateClient(gomock.Any(), gomock.Eq(arg)).
							  Times(1).
							  Return(client, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchClient(t, client, recorder.Body)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
	
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshall body data into JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/clients"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})	


	}
}

func TestGetClientAPI(t *testing.T) {
	client := randomClient()
	testCases := []struct{
		name string
		clientID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			clientID: client.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetClient(gomock.Any(), gomock.Eq(client.ID)).Times(1).Return(client, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchClient(t, client, recorder.Body)
			},
		},
		{
			name: "BadRequest",
			clientID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetClient(gomock.Any(), gomock.Any()).Times(0).Return(db.Client{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			clientID: client.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetClient(gomock.Any(), gomock.Eq(client.ID)).Times(1).Return(db.Client{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			clientID: client.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetClient(gomock.Any(), gomock.Eq(client.ID)).Times(1).Return(db.Client{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func (t *testing.T)  {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			url := fmt.Sprintf("/clients/%d", tc.clientID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateClientAPI(t *testing.T) {
	client := randomClient()
	testCases := []struct{
		name string
		body gin.H
		buildStubs func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			body: gin.H{
				"id": client.ID,
				"first_name": client.FirstName,
				"last_name": client.LastName,
				"country_id": client.CountryID,
				"active": client.Active,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.UpdateClientParams {
					ID: client.ID,
					FirstName: client.FirstName,
					LastName: client.LastName,
					CountryID: client.CountryID,
					Active: client.Active,
				}
				mockStore.EXPECT().UpdateClient(gomock.Any(), gomock.Eq(arg)).Times(1).Return(client, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchClient(t, client, recorder.Body)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, "/clients", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomClient() db.Client {
	return db.Client {
		ID: util.RandomInt(1, 1000),
		FirstName: util.RandomString(10),
		LastName: util.RandomString(10),
		CountryID: util.RandomInt(1, 228),
		UserID: util.RandomInt(1, 10),
	}
}

func requireBodyMatchClient(t *testing.T, client db.Client, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotClient db.Client
	err = json.Unmarshal(data, &gotClient)
	require.NoError(t, err)
	require.Equal(t, client, gotClient)
}