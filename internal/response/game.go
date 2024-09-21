package response

import (
	"time"

	"galere.se/oss-codenames-api/internal/domain/game"
)

//
// Types
//

type GameRoomResponse struct {
	Id           string                  `json:"id"`
	CreatedAt    time.Time               `json:"createdAt"`
	State        string                  `json:"state"`
	Name         string                  `json:"name"`
	RedTeam      []*PublicPlayerResponse `json:"redTeam"`
	BlueTeam     []*PublicPlayerResponse `json:"blueTeam"`
	Players      []*PublicPlayerResponse `json:"players"`
	CurrentRound *GameRoundResponse      `json:"currentRound"`
	RoundHistory []*GameRoundResponse    `json:"roundHistory"`
	RoomCode     int                     `json:"roomCode"`
}

type GameRoundResponse struct {
	RedScore       int                        `json:"redScore"`
	BlueScore      int                        `json:"blueScore"`
	CurrentTurn    *GameTurnResponse          `json:"currentTurn"`
	StartedAt      time.Time                  `json:"startedAt"`
	FinishedAt     time.Time                  `json:"finishedAt"`
	WinnerTeam     string                     `json:"winnerTeam"`
	RedSpymaster   *PublicPlayerResponse      `json:"redSpymaster"`
	BlueSpymaster  *PublicPlayerResponse      `json:"blueSpymaster"`
	RedScoreLimit  int                        `json:"redScoreLimit"`
	BlueScoreLimit int                        `json:"blueScoreLimit"`
	BoardTiles     map[int]*BoardTileResponse `json:"boardTiles"`
	GuessedTiles   map[int]bool               `json:"guessedTiles"`
	TurnHistory    []*GameTurnResponse        `json:"turnHistory"`
	WinReason      string                     `json:"winReason"`
}

type GameTurnResponse struct {
	CurrentTeam      string                  `json:"currentTeam"`
	StartedAt        time.Time               `json:"startedAt"`
	FinishedAt       time.Time               `json:"finishedAt"`
	Clue             string                  `json:"clue"`
	GuessAmount      int                     `json:"guessAmount"`
	GuessesMade      int                     `json:"guessesMade"`
	GuessedTiles     map[int]bool            `json:"guessedTiles"`
	PointedTiles     map[int]map[string]bool `json:"pointedTiles"`
	UnlimitedGuesses bool                    `json:"unlimitedGuesses"`
}

type BoardTileResponse struct {
	ImageUrl string `json:"imageUrl"`
}

//
// Instance builders
//

func NewGameRoomResponse(gameRoom *game.GameRoom) *GameRoomResponse {

	if gameRoom == nil {
		return nil
	}

	response := &GameRoomResponse{
		Id:           gameRoom.Id,
		CreatedAt:    gameRoom.CreatedAt,
		State:        string(gameRoom.State),
		Name:         gameRoom.Name,
		RedTeam:      mapPublicPlayersResponse(gameRoom.RedTeam),
		BlueTeam:     mapPublicPlayersResponse(gameRoom.BlueTeam),
		Players:      mapPublicPlayersResponse(gameRoom.Players),
		CurrentRound: NewGameRoundResponse(gameRoom.CurrentRound),
		RoundHistory: mapGameRoundsResponse(gameRoom.RoundHistory),
		RoomCode:     gameRoom.RoomCode,
	}

	return response
}

func NewGameRoundResponse(gameRound *game.GameRound) *GameRoundResponse {

	if gameRound == nil {
		return nil
	}

	return &GameRoundResponse{
		RedScore:       int(gameRound.RedScore),
		BlueScore:      int(gameRound.BlueScore),
		CurrentTurn:    NewGameTurnResponse(gameRound.CurrentTurn),
		StartedAt:      gameRound.StartedAt,
		FinishedAt:     gameRound.FinishedAt,
		WinnerTeam:     string(gameRound.WinnerTeam),
		RedSpymaster:   NewPublicPlayerResponse(gameRound.RedSpymaster),
		BlueSpymaster:  NewPublicPlayerResponse(gameRound.BlueSpymaster),
		RedScoreLimit:  int(gameRound.RedScoreLimit),
		BlueScoreLimit: int(gameRound.BlueScoreLimit),
		BoardTiles:     mapBoardTilesResponse(gameRound.BoardTiles),
		GuessedTiles:   gameRound.GuessedTiles,
		TurnHistory:    mapGameTurnsResponse(gameRound.TurnHistory),
		WinReason:      string(gameRound.WinReason),
	}
}

func NewGameTurnResponse(gameTurn *game.GameTurn) *GameTurnResponse {

	if gameTurn == nil {
		return nil
	}

	return &GameTurnResponse{
		CurrentTeam:      string(gameTurn.CurrentTeam),
		StartedAt:        gameTurn.StartedAt,
		FinishedAt:       gameTurn.FinishedAt,
		Clue:             gameTurn.Clue,
		GuessAmount:      gameTurn.GuessAmount,
		GuessesMade:      gameTurn.GuessesMade,
		GuessedTiles:     gameTurn.GuessedTiles,
		PointedTiles:     gameTurn.PointedTiles,
		UnlimitedGuesses: gameTurn.UnlimitedGuesses,
	}
}

func NewBoardTileResponse(boardTile *game.BoardTile) *BoardTileResponse {

	if boardTile == nil {
		return nil
	}

	return &BoardTileResponse{
		ImageUrl: boardTile.ImageUrl,
	}
}

//
// Mappers
//

func mapPublicPlayersResponse(players []*game.Player) []*PublicPlayerResponse {
	playerResponse := make([]*PublicPlayerResponse, len(players))

	for i, player := range players {
		playerResponse[i] = NewPublicPlayerResponse(player)
	}

	return playerResponse
}

func mapGameRoundsResponse(gameRounds []*game.GameRound) []*GameRoundResponse {

	gameRoundResponse := make([]*GameRoundResponse, len(gameRounds))

	for i, gameRound := range gameRounds {
		gameRoundResponse[i] = NewGameRoundResponse(gameRound)
	}

	return gameRoundResponse
}

func mapGameTurnsResponse(turnHistory []*game.GameTurn) []*GameTurnResponse {

	turnHistoryResponse := make([]*GameTurnResponse, len(turnHistory))

	for i, turn := range turnHistory {
		turnHistoryResponse[i] = NewGameTurnResponse(turn)
	}

	return turnHistoryResponse
}

func mapBoardTilesResponse(boardTiles map[int]game.BoardTile) map[int]*BoardTileResponse {

	boardTilesResponse := make(map[int]*BoardTileResponse)

	for tileId, tile := range boardTiles {
		boardTilesResponse[tileId] = NewBoardTileResponse(&tile)
	}

	return boardTilesResponse
}
