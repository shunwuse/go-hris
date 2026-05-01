package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode_UintValue(t *testing.T) {
	type testCursor struct {
		ID uint `json:"id"`
	}

	cursor := testCursor{ID: 123}
	encoded, err := Encode(cursor)

	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestEncode_StringValue(t *testing.T) {
	type testCursor struct {
		Token string `json:"token"`
	}

	cursor := testCursor{Token: "abc123"}
	encoded, err := Encode(cursor)

	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestDecode_ValidCursor(t *testing.T) {
	type testCursor struct {
		ID uint `json:"id"`
	}

	original := testCursor{ID: 456}
	encoded, err := Encode(original)
	require.NoError(t, err)

	decoded, err := Decode[testCursor](encoded)
	require.NoError(t, err)
	require.NotNil(t, decoded)
	assert.Equal(t, original.ID, decoded.ID)
}

func TestDecode_EmptyCursor(t *testing.T) {
	decoded, err := Decode[struct{}]("")
	assert.NoError(t, err)
	assert.Nil(t, decoded)
}

func TestDecode_InvalidBase64(t *testing.T) {
	decoded, err := Decode[struct{}]("invalid!!!base64")
	assert.Error(t, err)
	assert.Nil(t, decoded)
	assert.Contains(t, err.Error(), "failed to decode cursor")
}

func TestDecode_InvalidJSON(t *testing.T) {
	// Create a valid base64 string that's not valid JSON
	encoded := "aW52YWxpZCBqc29u" // "invalid json" in base64
	decoded, err := Decode[struct{ ID uint }](encoded)
	assert.Error(t, err)
	assert.Nil(t, decoded)
	assert.Contains(t, err.Error(), "failed to unmarshal cursor")
}

func TestEncodeDecode_RoundTrip(t *testing.T) {
	type testCursor struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
		Count int    `json:"count"`
	}

	original := testCursor{
		ID:    789,
		Token: "test-token",
		Count: 42,
	}

	// Encode
	encoded, err := Encode(original)
	require.NoError(t, err)

	// Decode
	decoded, err := Decode[testCursor](encoded)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.Token, decoded.Token)
	assert.Equal(t, original.Count, decoded.Count)
}

func TestEncodeDecode_WithNullFields(t *testing.T) {
	type testCursor struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}

	original := testCursor{ID: 100, Status: ""}

	encoded, err := Encode(original)
	require.NoError(t, err)

	decoded, err := Decode[testCursor](encoded)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.Status, decoded.Status)
}

func TestEncodeDecode_ConsistentAcrossMultipleCalls(t *testing.T) {
	type testCursor struct {
		ID uint `json:"id"`
	}

	original := testCursor{ID: 999}

	// Encode multiple times
	encoded1, err := Encode(original)
	require.NoError(t, err)

	encoded2, err := Encode(original)
	require.NoError(t, err)

	// Should be identical
	assert.Equal(t, encoded1, encoded2)

	// Decode both
	decoded1, err := Decode[testCursor](encoded1)
	require.NoError(t, err)

	decoded2, err := Decode[testCursor](encoded2)
	require.NoError(t, err)

	assert.Equal(t, decoded1, decoded2)
}
