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

			server := newTestServer(t, store)
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

			server := newTestServer(t, store)
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
		{
			name: "NotFound",
			body: gin.H{
				"id": 2000,
				"first_name": client.FirstName,
				"last_name": client.LastName,
				"country_id": client.CountryID,
				"active": client.Active,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.UpdateClientParams {
					ID: 2000,
					FirstName: client.FirstName,
					LastName: client.LastName,
					CountryID: client.CountryID,
					Active: client.Active,
				}
				mockStore.EXPECT().UpdateClient(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Client{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)				
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"id": 0,
				"first_name": client.FirstName,
				"last_name": client.LastName,
				"country_id": client.CountryID,
				"active": client.Active,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.UpdateClientParams {
					ID: 0,
					FirstName: client.FirstName,
					LastName: client.LastName,
					CountryID: client.CountryID,
					Active: client.Active,
				}
				mockStore.EXPECT().UpdateClient(gomock.Any(), gomock.Eq(arg)).Times(0).Return(client, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)				
			},
		},
		{
			name: "InternalServerError",
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
				mockStore.EXPECT().UpdateClient(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Client{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := newTestServer(t, mockStore)
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

type Query struct {
	Page int32
	PageSize int32
}

func TestListClientsAPI(t *testing.T) {
	randomCount := util.RandomInt(5, 10)
	clients := make([]db.Client, int(randomCount))
	
	for i := int64(0); i < randomCount; i++ {
		clients[i] = randomClient()
	}

	testCases := []struct{
		name string
		query Query
		buildStubs func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Page: 1,
				PageSize: int32(randomCount),
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.ListClientsParams {
					Limit: int32(randomCount),
					Offset: 0,
				}

				mockStore.EXPECT().ListClients(gomock.Any(), gomock.Eq(arg)).Times(1).Return(clients, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchClients(t, clients, recorder.Body)
			},
		},
		{
			name: "BadRequest",
			query: Query{
				Page: 0,
				PageSize: 0,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.ListClientsParams {
					Limit: 0,
					Offset: 0,
				}

				mockStore.EXPECT().ListClients(gomock.Any(), gomock.Eq(arg)).Times(0).Return([]db.Client{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			query: Query{
				Page: 1000,
				PageSize: int32(randomCount),
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.ListClientsParams {
					Limit: int32(randomCount),
					Offset: 999 * int32(randomCount),
				}

				mockStore.EXPECT().ListClients(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Client{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				Page: 1,
				PageSize: int32(randomCount),
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.ListClientsParams {
					Limit: int32(randomCount),
					Offset: 0,
				}

				mockStore.EXPECT().ListClients(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Client{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/clients", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.Page))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))

			req.URL.RawQuery = q.Encode()

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

func requireBodyMatchClients(t *testing.T, clients []db.Client, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotClients []db.Client
	err = json.Unmarshal(data, &gotClients)
	require.NoError(t, err)
	require.Equal(t, clients, gotClients)
}