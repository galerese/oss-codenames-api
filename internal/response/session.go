package response

import (
	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
)

//
// Types
//

type SessionResponse struct {
	Player      *PrivatePlayerResponse `json:"player"`
	CurrentRoom *GameRoomResponse      `json:"currentRoom"`
}

type PrivatePlayerResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
}

type PublicPlayerResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//
// Instance builders :)
//

func NewSessionResponse(session *session.Session) *SessionResponse {
	response := &SessionResponse{
		Player:      NewPrivatePlayerResponse(session.Player),
		CurrentRoom: NewGameRoomResponse(session.CurrentRoom),
	}

	return response
}

func NewPublicPlayerResponse(player *game.Player) *PublicPlayerResponse {

	if player == nil {
		return nil
	}

	return &PublicPlayerResponse{
		Id:   player.Id,
		Name: player.Name,
	}
}

func NewPrivatePlayerResponse(player *game.Player) *PrivatePlayerResponse {

	if player == nil {
		return nil
	}

	return &PrivatePlayerResponse{
		Id:    player.Id,
		Token: player.Token,
		Name:  player.Name,
	}
}
