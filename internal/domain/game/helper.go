package game

import (
	"fmt"
	"math/rand"

	"github.com/sirupsen/logrus"
)

var (
	// Borrowed from https://github.com/goombaio/namegenerator/blob/master/data.go
	// ADJECTIVES ...
	RandomGameRoomWordsAdjective = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer",
		"icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient",
		"twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing",
		"broken", "cold", "damp", "falling", "frosty", "long", "late", "lingering",
		"bold", "little", "morning", "muddy", "old", "rough", "still", "small",
		"sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black",
		"young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral",
		"restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}

	// NOUNS ...
	RandomGameRoomWordsNouns = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning",
		"snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest",
		"hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly",
		"bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass",
		"haze", "mountain", "night", "pond", "darkness", "snowflake", "silence",
		"sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower",
		"wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog",
		"frost", "voice", "paper", "frog", "smoke", "star"}
)

func (s *Service) generateRandomGameRoomName() (string, error) {

	// Randomly select two adjectives and a noun and the adjectives are not repeated
	adjective1Index := rand.Intn(len(RandomGameRoomWordsAdjective))
	adjective2Index := rand.Intn(len(RandomGameRoomWordsAdjective))
	for adjective1Index == adjective2Index {
		adjective2Index = rand.Intn(len(RandomGameRoomWordsAdjective))
	}

	adjective1 := RandomGameRoomWordsAdjective[adjective1Index]
	adjective2 := RandomGameRoomWordsAdjective[adjective2Index]
	noun := RandomGameRoomWordsNouns[rand.Intn(len(RandomGameRoomWordsNouns))]

	return fmt.Sprintf("%s-%s-%s", adjective1, adjective2, noun), nil
}

// placeholder for triggering events
func (s *Service) triggerGameRoomEvents(room *GameRoom, eventType GameRoomEvent) error {
	logrus.Debugf("Triggering event [%s] for room [%s]", eventType, room.Id)
	return nil
}

// Makes sure we have consistent room state :)
func (s *Service) ensureGameRoundAndTurnExist(room *GameRoom) error {
	if room.CurrentRound == nil {
		return NewUnexpectedError(nil, "Expected a game round to be created already!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return NewUnexpectedError(nil, "Expected a game turn to be created already!")
	}

	return nil
}
