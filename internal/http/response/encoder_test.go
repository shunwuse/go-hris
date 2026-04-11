package response

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	JSON(rr, http.StatusAccepted, map[string]string{"message": "ok"})

	require.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var body map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ok", body["message"])
}

func TestCreated_WithNilPayload(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	Created(rr, nil)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestNoContent(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	NoContent(rr)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestOffsetList(t *testing.T) {
	t.Parallel()

	meta := OffsetPaginationMeta{
		Total:       100,
		PerPage:     10,
		CurrentPage: 2,
		LastPage:    10,
	}

	rr := httptest.NewRecorder()
	OffsetList(rr, []string{"a", "b"}, meta)

	require.Equal(t, http.StatusOK, rr.Code)

	var body struct {
		Data []string             `json:"data"`
		Meta OffsetPaginationMeta `json:"meta"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, body.Data)
	assert.Equal(t, meta, body.Meta)
}

func TestCursorList(t *testing.T) {
	t.Parallel()

	meta := CursorPaginationMeta{
		NextCursor: "next-1",
		HasMore:    true,
	}

	rr := httptest.NewRecorder()
	CursorList(rr, []int{1, 2}, meta)

	require.Equal(t, http.StatusOK, rr.Code)

	var body struct {
		Data []int                `json:"data"`
		Meta CursorPaginationMeta `json:"meta"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, body.Data)
	assert.Equal(t, meta, body.Meta)
}

func TestError_DomainError(t *testing.T) {
	t.Parallel()

	errValidation := errors.ErrValidationFailed.WithDetails(map[string]string{
		"email": "required",
	})

	rr := httptest.NewRecorder()
	Error(rr, errValidation)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var body ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, errors.CodeValidationFailed, body.Error.Code)
	assert.Equal(t, "validation failed", body.Error.Message)
	assert.Equal(t, "required", body.Error.Details["email"])
}

func TestError_UnknownDomainCodeFallsBackToInternal(t *testing.T) {
	t.Parallel()

	appErr := errors.New("UNKNOWN_CODE", "unexpected")

	rr := httptest.NewRecorder()
	Error(rr, appErr)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var body ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "UNKNOWN_CODE", body.Error.Code)
	assert.Equal(t, "unexpected", body.Error.Message)
}

func TestError_UnexpectedError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	Error(rr, stderrors.New("boom"))

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var body ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, errors.CodeInternalError, body.Error.Code)
	assert.Equal(t, "boom", body.Error.Message)
}
