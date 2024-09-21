import { sleep, check } from 'k6';
import { Options } from 'k6/options';
import http from 'k6/http';
import { API } from './lib/api';

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

  const api = new API('http://localhost:3033');

  // We define game Id based on the current VU, because it's not possible to share states between VUs
  const gameId = Math.ceil(__VU / GAME_PLAYER_AMOUNT)
  const currentPlayerId = __VU % GAME_PLAYER_AMOUNT
  const currentPlayer = PLAYER_CONFIG[currentPlayerId]


  // Start by creating a new session
  const res = api.post('/v1/sessions');
  check(res, {
    'status is 200': () => res.status === 200,
  });

  sleep(1);
};
