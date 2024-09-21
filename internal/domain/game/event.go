package game

type GameRoomEvent string

const (
	GameRoomEventGameStarted       GameRoomEvent = "GAME_STARTED"
	GameRoomEventRoundStarted      GameRoomEvent = "ROUND_STARTED"
	GameRoomEventSpymasterSelected GameRoomEvent = "SPYMASTER_SELECTED"
	GameRoomEventSpymastersSettled GameRoomEvent = "SPYMASTERS_SETTLED"
	GameRoomEventClueSelected      GameRoomEvent = "CLUE_SELECTED"
	GameRoomEventRoundEnded        GameRoomEvent = "ROUND_ENDED"
	GameRoomEventGuessingStopped   GameRoomEvent = "GUESSING_STOPPED"
	GameRoomEventTeamSelected      GameRoomEvent = "TEAM_SELECTED"
	GameRoomEventPlayerJoined      GameRoomEvent = "PLAYER_JOINED"
	GameRoomEventTurnEnded         GameRoomEvent = "TURN_ENDED"

	GameRoomEventTilePointed   GameRoomEvent = "TILE_POINTED"
	GameRoomEventTileUnpointed GameRoomEvent = "TILE_UNPOINTED"
	GameRoomEventTileGuessed   GameRoomEvent = "TILE_GUESSED"
)
