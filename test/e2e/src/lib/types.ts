// Create types based on the following golang API response types

export type TeamName = 'red' | 'blue'

export type Session = {
    player: Player
    currentRoom: GameRoom
}

export type GameRoom = {
    id: string
    createdAt: Date
    state: string
    name: string
    redTeam: Player[]
    blueTeam: Player[]
    players: Player[]
    currentRound: GameRound
    roundHistory: GameRound[]
}

export type GameRound = {
    redScore: number
    blueScore: number
    currentTurn: GameTurn
    startedAt: Date
    finishedAt: Date
    winnerTeam: string
    redSpymaster: Player
    blueSpymaster: Player
    redScoreLimit: number
    blueScoreLimit: number
    boardTiles: BoardTile[]
    guessedTiles: boolean[]
    turnHistory: GameTurn[]
    winReason: string
}

export type GameTurn = {
    currentTeam: string
    startedAt: Date
    finishedAt: Date
    clue: string
    guessAmount: number
    guessesMade: number
    guessedTiles: boolean[]
    pointedTiles: boolean[]
    unlimitedGuesses: boolean
}

export type BoardTile = {
    imageUrl: string
}

export type Player = {
    id: string
    name: string
    token?: string
}

    

