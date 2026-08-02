package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddd/internal/application/dto"
	"goddd/internal/domain/models"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Failed(t *testing.T) {
	req := dto.CreateUserRequest{Username: "", Email: gofakeit.Email()}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users"), "application/json", reader)

	// assertion
	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var data map[string]string
	require.NoError(t, json.Unmarshal(resData, &data))

	expRes := map[string]string{
		"code":    "INVALID_INPUT",
		"message": "invalid input: username is required",
	}
	assert.Equal(t, expRes, data)
}

func TestCreateUser_Succeed(t *testing.T) {
	req := dto.CreateUserRequest{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
	}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users"), "application/json", reader)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	// assertion
	resPostData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resPostUser models.User
	require.NoError(t, json.Unmarshal(resPostData, &resPostUser))

	assert.NotEqual(t, uuid.Nil, resPostUser.ID)
	assert.Equal(t, req.Username, resPostUser.Username)
	assert.Equal(t, req.Email, resPostUser.Email)
	assert.False(t, resPostUser.CreatedAt.IsZero())

	// try to retrieve user
	res, err = http.Get(fmt.Sprintf("%s/users/%s", e2eBaseUrl, resPostUser.ID.String()))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

	// assertion
	resGetData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data map[string]string
	require.NoError(t, json.Unmarshal(resGetData, &data))

	expRes := map[string]string{
		"ID":        resPostUser.ID.String(),
		"Email":     req.Email,
		"Username":  req.Username,
		"CreatedAt": resPostUser.CreatedAt.Format("2006-01-02T15:04:05.999999Z"),
		"UpdatedAt": resPostUser.UpdatedAt.Format("2006-01-02T15:04:05.999999Z"),
	}
	assert.Equal(t, expRes, data)
}

func TestBulkCreateUsers_Failed(t *testing.T) {
	req := dto.CreateUsersRequest{{
		Username: gofakeit.Username(),
		Email:    "",
	}}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users/bulk"), "application/json", reader)

	// assertion
	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var data map[string]string
	require.NoError(t, json.Unmarshal(resData, &data))

	expRes := map[string]string{
		"code":    "INVALID_INPUT",
		"message": "invalid input: email is not a valid address: email=''",
	}
	assert.Equal(t, expRes, data)
}

func TestBulkCreateUsers_Succeed(t *testing.T) {
	req := dto.CreateUsersRequest{{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
	}}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users/bulk"), "application/json", reader)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	// assertion
	resPostData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resPostUsers []models.User
	require.NoError(t, json.Unmarshal(resPostData, &resPostUsers))

	require.Len(t, resPostUsers, 1)
	resPostUser := resPostUsers[0]
	assert.NotEqual(t, uuid.Nil, resPostUser.ID)
	assert.Equal(t, req[0].Username, resPostUser.Username)
	assert.Equal(t, req[0].Email, resPostUser.Email)
	assert.False(t, resPostUser.CreatedAt.IsZero())
}

func TestFailedToGetUser(t *testing.T) {
	res, err := http.Get(fmt.Sprintf(e2eBaseUrl + "/users/invalid-id"))
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	// assertion
	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data map[string]string
	require.NoError(t, json.Unmarshal(resData, &data))

	expRes := map[string]string{
		"code":    "INVALID_INPUT",
		"message": "invalid input: user Id must be of type uuid",
	}
	assert.Equal(t, expRes, data)
}

func TestListUsers(t *testing.T) {
	req := dto.CreateUserRequest{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
	}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users"), "application/json", reader)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	// assertion
	resPostData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resPostUser models.User
	require.NoError(t, json.Unmarshal(resPostData, &resPostUser))

	// try to list users
	res, err = http.Get(fmt.Sprintf("%s/users/?limit=1", e2eBaseUrl))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

	// assertion
	resGetData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var list dto.ListUsersResponse
	require.NoError(t, json.Unmarshal(resGetData, &list))

	assert.Equal(t, dto.ListUsersResponse{
		Users:      []models.User{resPostUser},
		Limit:      1,
		NextCursor: lo.ToPtr(resPostUser.ID.String()),
	}, list)
}

func TestUpdateUser_Failed(t *testing.T) {
	payload, err := json.Marshal(dto.UpdateUserRequest{})
	require.NoError(t, err)

	reader := bytes.NewReader(payload)
	res, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users/invalid-id"), "application/json", reader)

	// assertion
	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var data map[string]string
	require.NoError(t, json.Unmarshal(resData, &data))

	expRes := map[string]string{
		"code":    "INVALID_INPUT",
		"message": "invalid input: user Id must be of type uuid",
	}
	assert.Equal(t, expRes, data)
}

func TestUpdateUser_Succeed(t *testing.T) {

	// --- 1. create a new user ---
	createReq := dto.CreateUserRequest{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
	}

	createPayload, err := json.Marshal(createReq)
	require.NoError(t, err)

	reader := bytes.NewReader(createPayload)
	createRes, err := http.Post(fmt.Sprintf(e2eBaseUrl+"/users"), "application/json", reader)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, createRes.StatusCode)

	// assertion
	createResPostData, err := io.ReadAll(createRes.Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", createRes.Header.Get("Content-Type"))

	var createResPostUser models.User
	require.NoError(t, json.Unmarshal(createResPostData, &createResPostUser))

	// --- 2. update the created user ---
	updateReq := dto.UpdateUserRequest{Username: "john"}
	updatePayload, err := json.Marshal(updateReq)
	require.NoError(t, err)

	reader = bytes.NewReader(updatePayload)
	updateRes, err := http.Post(fmt.Sprintf("%s/users/%s", e2eBaseUrl, createResPostUser.ID.String()), "application/json", reader)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, updateRes.StatusCode)

	// assertion
	updateResPostData, err := io.ReadAll(updateRes.Body)
	require.NoError(t, err)
	require.Equal(t, "application/json", updateRes.Header.Get("Content-Type"))

	var updateResPostUser models.User
	require.NoError(t, json.Unmarshal(updateResPostData, &updateResPostUser))
}
