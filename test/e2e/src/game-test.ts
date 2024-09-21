import { Options } from 'k6/options';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';



export let options:Options = {
  vus: 50,
  duration: '10s'
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

// Tests a single game
export default () => {

  const api = new Httpx({ baseURL: 'http://localhost:3033' });

  // We define game Id based on the current VU, because it's not possible to share states between VUs
  const gameId = Math.ceil(__VU / GAME_PLAYER_AMOUNT)
  const currentPlayerId = ((__VU - 1) % GAME_PLAYER_AMOUNT) + 1
  const currentPlayer = PLAYER_CONFIG[currentPlayerId]

  const debug = function(...args: any[]){
    console.log.apply(console, [`VU:${__VU} gameId=${gameId} playerId=${currentPlayerId}`, ...args])
  }

  const context: any = {
    session: null,
  }

  // Start by creating a new session

  describe('Session creation', () => {
    const res = api.post('/v1/sessions');
    expect(res.status, 'response status').to.equal(200)

    const body = JSON.parse(res.body as string)

    expect(body, 'body').to.be.a('object')
    debug('Session creation', body)

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

  // PLAYER A
  // Room Creation :)
  if (currentPlayerId === 1) {

    describe('Room creation with no player name', () => {
      const res = api.post('/v1/rooms');
      expect(res.status, 'response status').to.equal(400)
    });

    describe('Room creation', () => {
      const res = api.post('/v1/rooms', {
        playerName: currentPlayer.name
      });

      const body = JSON.parse(res.body as string)
  
      expect(body, 'body').to.be.a('object')
      debug('Room creation', res.body)

      expect(res.status, 'response status').to.equal(201)
    });
  }


  
};
