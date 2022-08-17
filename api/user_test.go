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
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher { 
	return eqCreateUserParamsMatcher{arg, password} 
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	testCases := []struct{
		name string		
		body gin.H
		buildStubs func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"password_repeated": user.Password,
				"email": user.Email,
			},			
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams {
					Username: user.Username,					
					Email: user.Email,
				}
				mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, user.Password)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserBodyMatch(t, user, recorder.Body)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"username": "",
				"password": user.Password,
				"password_repeated": user.Password,
				"email": user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams {
					Username: user.Username,					
					Email: user.Email,
				}
				mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, user.Password)).Times(0).Return(db.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotUnique",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"password_repeated": user.Password,
				"email": user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams {
					Username: user.Username,					
					Email: user.Email,
				}
				mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, user.Password)).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "#",
				"password": user.Password,
				"password_repeated": user.Password,
				"email": user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"password_repeated": user.Password,
				"email": "something",
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
				"password_repeated": user.Password,
				"email": user.Email,
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateUserParams {
					Username: user.Username,					
					Email: user.Email,
				}
				mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, user.Password)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, req)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireUserBodyMatch(t *testing.T, user db.User, buffer *bytes.Buffer) {
	data, err := ioutil.ReadAll(buffer)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.False(t, gotUser.CreatedAt.IsZero())
}

func createRandomUser(t *testing.T) db.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	return db.User{
		ID: util.RandomInt(1, 100),
		Username: util.RandomString(6),
		Password: hashedPassword,
		Email: util.RandomEmail(6),
		CreatedAt: time.Now(),
	}
}