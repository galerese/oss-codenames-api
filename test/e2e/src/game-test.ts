import { Options } from 'k6/options';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';
import redis from 'k6/experimental/redis';
import { validateEmptyRoom } from './lib/validation';
import { GameService } from './lib/game-service';

// TODO
// - validate that two rooms have different names :)
// - validate that a player cannot create a room if he already has one

export let options: Options = {
  vus: 4,
  duration: '1s',
  iterations: 4
};

const GAME_PLAYER_AMOUNT = 4
const PLAYER_CONFIG = [
  {
    name: 'Player A'
  },
  {
    name: 'Player B'
  },
  {
    name: 'Player C'
  },
  {
    name: 'Player D'
  },
]

type sharedData = {
  redisHash: string
}

export function setup(): sharedData {

  return {
    redisHash: Math.random().toString(36).substring(2, 8)
  }
}


// Tests a single game
export default async function (data: sharedData) {

  const redisClient = new redis.Client('redis://localhost:6379');

  const api = new Httpx({ baseURL: 'http://localhost:3033' });
  api.addHeader('Content-Type', 'application/json')

  const gameService = new GameService(api)

  // We define game Id based on the current VU, because it's not possible to share states between VUs
  const gameId = Math.ceil(__VU / GAME_PLAYER_AMOUNT)
  const currentPlayerId = ((__VU - 1) % GAME_PLAYER_AMOUNT) + 1
  const currentPlayer = PLAYER_CONFIG[currentPlayerId - 1]

  const REDIS_KEY_GAME_ROOM_NAME = `${data.redisHash}-game-room-name-${gameId}`

  const debug = function (...args: any[]) {
    console.log.apply(console, [`VU:${__VU} gameId=${gameId} playerId=${currentPlayerId}`, ...args])
  }

  const debugResponse = function (msg: string, res: any) {
    debug(msg, res.body, res.status)
  }

  const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

  const retryWithTimeout = async (fn: () => Promise<any>, timeout: number, err: string) => {
    const start = Date.now()
    while (Date.now() - start < timeout) {
      try {

        const res = await fn()
        if (res != undefined) {
          return res
        }
      } catch (e) {
        debug("Retrying with timeout.. error: ", e)
      }
      await delay(100)
    }
    throw new Error(err)
  }

  api.postJson = (path: string, body: any) => {
    return api.post(path, JSON.stringify(body))
  }

  api.patchJson = (path: string, body: any) => {
    return api.patch(path, JSON.stringify(body))
  }

  const context: any = {
    session: null,
    room: null
  }

  // Start by creating a new session

  describe('It should be possible to create a session', () => {
    const res = api.postJson('/v1/sessions')
    debugResponse('It should be possible to create a session', res)

    expect(res.status, 'response status').to.equal(200)

    const body = JSON.parse(res.body as string)

    expect(body, 'body').to.be.a('object')

    expect(body.player.id, 'player id').to.be.a('string')
    expect(body.player.token.length, 'player token id').to.be.greaterThan(0)

    expect(body.player.token, 'player token').to.be.a('string')
    expect(body.player.token.length, 'player token length').to.be.greaterThan(0)

    expect(body.player.token === body.player.id, 'token is not equal to id').to.be.false

    context.session = body
  });

  // Set client token :)
  expect(context.session?.player?.token, 'session token').to.be.not.null
  api.addHeader('Authorization', `Bearer ${context.session.player.token}`)


  // PLAYER A creates a room
  if (currentPlayerId === 1) {

    describe('It should not be possible to create a room with no player name', () => {
      const res = api.post('/v1/rooms');
      expect(res.status, 'response status').to.equal(400)
    });

    describe('It should be possible to create a room if we don\'t have one yet', () => {
      const res = api.postJson('/v1/rooms', {
        playerName: currentPlayer.name
      });

      debugResponse('It should be possible to create a room if we don\'t have one yet', res)

      const room = JSON.parse(res.body as string)
      expect(room, 'body').to.be.a('object')

      expect(res.status, 'response status').to.equal(201)

      validateEmptyRoom(room)

      // Validate room players
      expect(room.players.length, 'room players length').to.be.equal(1)
      expect(room.players[0].id, 'room players id').to.be.a('string')
      expect(room.players[0].id.length, 'room players id length').to.be.greaterThan(0)
      expect(room.players[0].name, 'room players name').to.be.equal(currentPlayer.name)

      context.room = room
    });

    // Set room name for this game :)
    await redisClient.set(REDIS_KEY_GAME_ROOM_NAME, context.room.name, 0)

  }
  // Other players just join the room
  else {

    const roomName = await retryWithTimeout(async () => {
      return redisClient.get(`${REDIS_KEY_GAME_ROOM_NAME}`)
    }, 2000, 'Could not find a room to join within 2s')

    // Get room details and validate it's okay :)
    describe('It should not be possible to get room details before joining', () => {
      const res = api.get(`/v1/rooms/${roomName}`);
      debugResponse('It should not be possible to get room details before joining', res)

      expect(res.status, 'response status').to.equal(403)
    })

    describe('Join room', () => {

      const res = api.postJson(`/v1/rooms/${roomName}/players`, {
        playerName: currentPlayer.name
      });
      debugResponse('Join room', res)

      expect(res.status, 'response status').to.equal(200)

      const room = JSON.parse(res.body as string)
      validateEmptyRoom(room)

      expect(room.players.length, 'room players length').to.be.greaterThan(1)
      expect(!!room.players.find((p: any) => p.id === context.session.player.id), 'player should be in the room').to.be.true

      context.room = room
    });


    // Get room details by name 
    describe('It should be possible to get room details after joining :)', () => {
      const res = api.get(`/v1/rooms/${context.room.name}/`);
      debugResponse('It should be possible to get room details after joining :)', res)

      expect(res.status, 'response status').to.equal(200)

      const room = JSON.parse(res.body as string)
      validateEmptyRoom(room)

      expect(room.players.length, 'room players length').to.be.greaterThan(1)
      expect(room.players.find((p: any) => p.name === currentPlayer.name), 'player should be in the room').to.be.not.null
    })

  }


  // Validate room :)
  expect(context.room, 'room').to.be.not.null

  //
  // Wait for all players to join the room (TODO LISTEN TO EVENTS)
  //
  await delay(1000)

  //
  // All players validate their session
  //
  describe('Session should contain user and room', () => {
    const res = api.get('/v1/session');
    debugResponse('Session should contain user and room', res)

    expect(res.status, 'response status').to.equal(200)

    const body = JSON.parse(res.body as string)
    expect(body, 'body').to.be.a('object')

    expect(body.player?.id, 'player id').to.be.a('string')
    expect(body.player?.id.length, 'player id length').to.be.greaterThan(0)

    expect(body.player?.token, 'player token').to.be.a('string')
    expect(body.player?.token.length, 'player token length').to.be.greaterThan(0)

    expect(body.player?.name, 'player name').to.be.equal(currentPlayer.name)

    expect(body.currentRoom, 'current room').to.be.not.null
    expect(body.currentRoom.id, 'current room id').to.be.equal(context.room.id)

    expect(body.currentRoom.players.length, 'current room players length').to.be.equal(PLAYER_CONFIG.length)
    expect(body.currentRoom.redTeam.length, 'current room red team length').to.be.equal(0)
    expect(body.currentRoom.blueTeam.length, 'current room blue team length').to.be.equal(0)
  });

  // Validate invalid state transition
  describe('It should not be possible change the game state to a random state name', () => {
    const res = api.patchJson(`/v1/rooms/${context.room.name}`, {
      state: 'BLA'
    })
    debugResponse('It should not be possible to change the game state to a random state name', res)

    expect(res.status, 'response status').to.equal(400)
  })

  // Validate invalid state transition
  describe('It should not be possible to set spymasters settled when the game hasn´t started yet', () => {
    const res = api.patchJson(`/v1/rooms/${context.room.name}`, {
      state: 'SPYMASTERS_SETTLED'
    })
    debugResponse('It should not be possible to set spymasters settled when the game hasn´t started yet', res)

    expect(res.status, 'response status').to.equal(400)
  })

  // All users try to start the game
  describe('It should be possible to start the game', async () => {
    const res = await api.patchJson(`/v1/rooms/${context.room.name}`, {
      state: 'STARTED'
    })
    debugResponse('It should be possible to start the game', res)

    await redisClient.incr(`${REDIS_KEY_GAME_ROOM_NAME}-start-game-response-${res.status}`)
  })

  await delay(1000)
  
  // Validate status after all users have tried to start the game
  expect(await redisClient.get(`${REDIS_KEY_GAME_ROOM_NAME}-start-game-response-400`), 'start game response 400').to.be.equal(PLAYER_CONFIG.length - 1)
  expect(await redisClient.get(`${REDIS_KEY_GAME_ROOM_NAME}-start-game-response-200`), 'start game response 200').to.be.equal(1)



};
