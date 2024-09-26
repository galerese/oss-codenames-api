import { GameRoom, Session, TeamName } from "./types"

type Response<T> = {
    response: any
    json: T
}

export class GameService {

    constructor(private readonly client: any) {
        this.client = client
    }

    async post(path: string, body?: any) {
        return await this.client.post(path, JSON.stringify(body || {})).then((res: any) => ({ response: res, json: JSON.parse(res.body) }))
    }

    async patch(path: string, body?: any) {
        return await this.client.patch(path, JSON.stringify(body || {})).then((res: any) => ({ response: res, json: JSON.parse(res.body) }))
    }

    async get(path: string) {
        return await this.client.get(path).then((res: any) => JSON.parse(res.body))
    }

    async createSession(): Promise<Response<Session>> {
        return this.post('/v1/session')
    }

    async getSession(): Promise<Response<Session>> {
        return this.get('/v1/session')
    }

    //
    // Game setup
    //

    async createRoom(playerName: string): Promise<Response<GameRoom>> {
        return this.post('/v1/rooms', { playerName })
    }

    async joinRoom(roomId: string, playerName: string): Promise<Response<GameRoom>> {
        return this.post(`/v1/rooms/${roomId}/join`, { playerName })
    }

    async getRoom(roomId: string): Promise<Response<GameRoom>> {
        return this.get(`/v1/rooms/${roomId}`)
    }

    async setTeam(roomId: string, playerId: string, team: TeamName): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}/players/${playerId}`, { team: team })
    }

    async startGame(roomId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}`, { state: 'GAME_STARTED' })
    }

    //
    // Spymaster selection
    //

    async setSpymaster(roomId: string, playerId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}/players/${playerId}`, { spymaster: true })
    }

    async settleSpymasters(roomId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}`, { state: 'SPYMASTERS_SETTLED' })
    }

    //
    // Clue & Guessing phase
    //
    async selectClue(roomId: string, clue: string, guessAmount: number): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}`, { state: 'CLUE_SELECTED', guessAmount, clue })
    }

    async guessTile(roomId: string, tileId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}/tiles/${tileId}`, { guessed: true })
    }

    async pointTile(roomId: string, tileId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}/tiles/${tileId}`, { pointed: true })
    }

    async unpointTile(roomId: string, tileId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}/tiles/${tileId}`, { pointed: false })
    }

    async endGuesses(roomId: string): Promise<Response<GameRoom>> {
        return this.patch(`/v1/rooms/${roomId}`, { state: 'GUESSING_STOPPED' })
    }

}