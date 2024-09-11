package game

import "time"

type Score int
type Identity string // "red", "blue", "neutral", "assassin"
type TeamName string // "red", "blue"
type TileCoordinate int

type GameRoomState string

const (
	//
	GameRoomStateWaitingForPlayers = "WAITING_FOR_PLAYERS"
	GameRoomStatePlaying           = "SELECT_SPYMASTERS"
	GameRoomStateSelectClue        = "SELECT_CLUE"
	GameRoomStateSelectGuess       = "SELECT_GUESS"
	GameRoomStateRoundEnded        = "ROUND_ENDED"
	// Transitional states
	GameRoomStateStarted            = "STARTED"
	GameRoomStateSpymastersSelected = "SPYMASTERS_SELECTED"
	GameRoomStateClueSelected       = "CLUE_SELECTED"
)

type GameRoom struct {
	CreatedAt    time.Time
	State        GameRoomState
	Name         string
	RedTeam      []Player
	BlueTeam     []Player
	Players      []Player
	CurrentRound GameRound
	RoundHistory []GameRound
	RoomCode     int
}

type GameRound struct {
	RedScore       Score
	BlueScore      Score
	CurrentTurn    GameTurn
	StartedAt      time.Time
	FinishedAt     time.Time
	WinnerTeam     TeamName
	RedSpymaster   Player
	BlueSpymaster  Player
	RedScoreLimit  Score
	BlueScoreLimit Score
	Identities     []Identity
	BoardTiles     []BoardTile
	GuessedTiles   []bool
	TurnHistory    []GameTurn
	WinReason      string
}

type GameTurn struct {
	CurrentTeam  TeamName
	StartedAt    time.Time
	FinishedAt   time.Time
	Clue         string
	GuessAmount  int
	GuessedTiles []TileCoordinate
	PointedTiles [][]Player
}

type BoardTile struct {
	ImageUrl string
}

type Player struct {
	Name  string
	IP    string
	Id    string
	Token string
}
