package game

const (
	DefaultBoardTileCount = 25
	MaxClueLength         = 20
	MaxGuessAmount        = 10
)

type WinReason string

const (
	WinReasonScore    = "SCORE"
	WinReasonAssassin = "ASSASSIN"
)

type Score int

type Identity string // "red", "blue", "neutral", "assassin"
const (
	IdentityRed      Identity = "red"
	IdentityBlue     Identity = "blue"
	IdentityNeutral  Identity = "neutral"
	IdentityAssassin Identity = "assassin"
)

type TeamName string // "red", "blue"
const (
	TeamNameRed  TeamName = "red"
	TeamNameBlue TeamName = "blue"
)

type TileCoordinate int

type GameRoomState string

const (
	//
	GameRoomStateWaitingForPlayers = "WAITING_FOR_PLAYERS"
	GameRoomStateSelectSpymasters  = "SELECT_SPYMASTERS"
	GameRoomStateSelectClue        = "SELECT_CLUE"
	GameRoomStateSelectGuess       = "SELECT_GUESS"
	GameRoomStateRoundEnded        = "ROUND_ENDED"
	// Transitional states
	GameRoomStateStarted           = "STARTED"
	GameRoomStateSpymastersSettled = "SPYMASTERS_SETTLED"
	GameRoomStateClueSelected      = "CLUE_SELECTED"
	GameRoomStateGuessingStopped   = "GUESSING_STOPPED"
)
