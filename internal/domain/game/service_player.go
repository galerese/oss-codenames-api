package game

import (
	"context"
	"fmt"

	"galere.se/oss-codenames-api/pkg/domain_util"
	"github.com/sirupsen/logrus"
)

// Adds an existing player to an existing game room
func (s *Service) AddPlayerToGameRoom(ctx context.Context, room *GameRoom, player *Player) (*GameRoom, error) {
	logrus.Infof("Adding player [%s] to game room [%s]", player.Name, room.Name)

	if room.State != GameRoomStateWaitingForPlayers {
		return nil, domain_util.NewStateValidationError("The game has already started, new players may not join now!")
	}

	// Nothing to do if the player is already in the room :)
	if room.IsPlayerInRoom(player) {
		logrus.Debugf("Player [%s] who is trying to join is already in game room [%s]", player.Id, room.Name)
		return room, nil
	}

	room.Players = append(room.Players, player)

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for new player")
	}

	logrus.Infof("Added player [%s] to game room [%s]", player.Name, room.Name)

	s.triggerGameRoomEvents(room, GameRoomEventPlayerJoined)

	return room, nil
}

// Updates the player team in a game room
func (s *Service) AddPlayerToTeam(ctx context.Context, room *GameRoom, player *Player, team TeamName) (*GameRoom, error) {
	logrus.Infof("Updating player [%s] team to [%s] in game room [%s]", player.Name, team, room.Name)

	if room.State != GameRoomStateWaitingForPlayers {
		return nil, domain_util.NewStateValidationError("The game has already started, you may not select your team at the moment!")
	}

	// Remove player from the other team if needed, or return early if the player is already in the desired team
	for i, p := range room.RedTeam {
		if p.Id == player.Id {
			if team == TeamNameRed {
				return room, nil
			}

			room.RedTeam = append(room.RedTeam[:i], room.RedTeam[i+1:]...)
			break
		}
	}

	for i, p := range room.BlueTeam {
		if p.Id == player.Id {
			if team == TeamNameBlue {
				return room, nil
			}

			room.BlueTeam = append(room.BlueTeam[:i], room.BlueTeam[i+1:]...)
			break
		}
	}

	// Otherwise, add the player to the desired team :)
	switch team {
	case TeamNameRed:
		room.RedTeam = append(room.RedTeam, player)
	case TeamNameBlue:
		room.BlueTeam = append(room.BlueTeam, player)
	}

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for updated player team")
	}

	logrus.Infof("Updated player [%s] team to [%s] in game room [%s]", player.Name, team, room.Name)

	s.triggerGameRoomEvents(room, GameRoomEventTeamSelected)

	return room, nil
}

// Sets the spymaster for the current game round
func (s *Service) SetSpymaster(ctx context.Context, room *GameRoom, player *Player) (*GameRoom, error) {
	logrus.Infof("Setting spymaster for game in room [%s] to player [%s]", room.Id, player.Name)

	// Validation

	if room.State != GameRoomStateSelectSpymasters {
		return nil, domain_util.NewStateValidationError("It's not currently time to select spymasters! You can only do this right before the game starts.")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before selecting spymasters!")
	}

	team := room.getPlayerTeam(player)
	if team == "" {
		return nil, domain_util.NewInvalidActionError("Provided player is not in the game room!")
	}

	// Setting the spymaster
	switch team {
	case TeamNameRed:
		if room.CurrentRound.RedSpymaster != nil && room.CurrentRound.RedSpymaster.Id != player.Id {
			logrus.Infof("Spymaster for red team already set to [%s], changing it", room.CurrentRound.RedSpymaster.Name)
		}
		room.CurrentRound.RedSpymaster = player

	case TeamNameBlue:
		if room.CurrentRound.BlueSpymaster != nil && room.CurrentRound.BlueSpymaster.Id != player.Id {
			logrus.Infof("Spymaster for blue team already set to [%s], changing it", room.CurrentRound.BlueSpymaster.Name)
		}
		room.CurrentRound.BlueSpymaster = player

	default:
		return nil, domain_util.NewInvalidParameterError(fmt.Sprintf("Unexpected team [%s] provided for spymaster selection!", team))
	}

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room when updating spymaster")
	}

	logrus.Infof("Set spymaster for game in room [%s] and player [%s] on team [%s]", room.Id, player.Name, team)

	s.triggerGameRoomEvents(room, GameRoomEventSpymasterSelected)

	return room, nil
}
