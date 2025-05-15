package rid_test

import (
	"strings"
	"testing"

	"github.com/onexstack/fastgo/internal/pkg/rid"
	"github.com/stretchr/testify/assert"
)

func Salt() string {
	return "staticSalt"
}

func TestResourceID_String(t *testing.T) {
	userID := rid.UserID
	assert.Equal(t, "user", userID.String(), "UserID.String() should return 'user'")

	postID := rid.PostID
	assert.Equal(t, "post", postID.String(), "PostID.String() should return 'post'")
}

func TestResourceID_New(t *testing.T) {
	userID := rid.UserID
	uniqueID := userID.New(1)
	assert.True(t, len(uniqueID) > 0, "Generated ID should not be empty")
	assert.Contains(t, uniqueID, "user-", "Generated ID should contain 'user-' prefix")

	anotherID := userID.New(2)
	assert.NotEqual(t, uniqueID, anotherID, "Generated IDs should be unique")
}

func BenchmarkResourceID_New(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := rid.UserID
		_ = userID.New(uint64(i))
	}
}

func FuzzResourceID_New(f *testing.F) {
	f.Add(uint64(1))
	f.Add(uint64(123456))
	f.Fuzz(func(t *testing.T, counter uint64) {
		result := rid.UserID.New(counter)
		assert.NotEmpty(t, result, "the generated unique identifier should not be empty")
		assert.Contains(t, result, rid.UserID.String()+"-", "the generated unique identifier should contain the prefix")
		splitParts := strings.SplitN(result, "-", 2)
		assert.Equal(t, rid.UserID.String(), splitParts[0], "the prefix of the generated unique identifier should match the expected prefix")
		if len(splitParts) == 2 {
			assert.Equal(t, 6, len(splitParts[1]), "the suffix of the generated unique identifier should be 6 characters long")
		} else {
			t.Error("The format of the generated unique identifier does not meet expectation")
		}
	})
}
