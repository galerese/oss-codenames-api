package game

import (
	"context"
	"strings"
	"testing"

	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/stretchr/testify/require"
)

type MockRepository struct {
}

func (r *MockRepository) GetGameRoomByCode(ctx context.Context, code string) (*GameRoom, error) {
	return nil, nil
}

func (r *MockRepository) GetGameRoomByName(ctx context.Context, name string) (*GameRoom, error) {
	return nil, nil
}

func (r *MockRepository) SaveGameRoom(ctx context.Context, gameRoom *GameRoom) error {
	return nil
}

func (r *MockRepository) GetRandomBoardTiles(ctx context.Context, count int) (map[int]BoardTile, error) {
	return nil, nil
}

func Test_generateRandomGameRoomName(t *testing.T) {

	assert := require.New(t)
	logger := logging.New("debug", "console")

	service := NewService(&MockRepository{}, &logger)

	names := map[string]bool{}

	// Run test multiple times to ensure randomness
	for i := 0; i < 1000; i++ {
		name, err := service.generateRandomGameRoomName()

		assert.NoError(err)

		parts := strings.Split(name, "-")
		assert.Len(parts, 3)
		assert.Contains(RandomGameRoomWordsAdjective, parts[0])
		assert.Contains(RandomGameRoomWordsAdjective, parts[1])
		assert.Contains(RandomGameRoomWordsNouns, parts[2])

		assert.NotEqual(parts[0], parts[1])

		names[name] = true
	}

	// Ensure names are unique to some extent
	assert.Greater(len(names), 900)

}
