package api

import (
	mockdb "bank-mvp/db/mock"
	db "bank-mvp/db/sqlc"
	"bank-mvp/util"
	"bytes"
	"encoding/json"
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
				requireBodyMatchClient(t, recorder.Body, client)
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

func randomClient() db.Client {
	return db.Client {
		FirstName: util.RandomString(10),
		LastName: util.RandomString(10),
		CountryID: util.RandomInt(1, 228),
		UserID: util.RandomInt(1, 10),
	}
}

func requireBodyMatchClient(t *testing.T, body *bytes.Buffer, client db.Client) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotClient db.Client
	err = json.Unmarshal(data, &gotClient)
	require.NoError(t, err)
	require.Equal(t, client, gotClient)
}