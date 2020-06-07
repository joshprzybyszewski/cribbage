import { push } from 'connected-react-router';
import { put, takeLatest } from 'redux-saga/effects';
import axios from 'axios';

import { game } from './types';
import { alertActions } from './actions';

// exit game
export function* watchGameExit() {
  yield takeLatest(game.EXIT_GAME, gameExit);
}
export function* gameExit() {
  yield put(push('/home')); // { type: game.reducer.EXIT_GAME });
}

// view game
export function* watchGameView() {
  yield takeLatest(game.VIEW_GAME, viewGame);
}
export function* viewGame({ payload }) {
  try {
    const res = yield axios.get(`/game/${payload}`);
    const gameID = res.data.GameID;
    yield put({
      type: game.reducer.VIEW_GAME,
      payload: { 
          gameID: gameID,
          gameJSON: res.data,
        },
    });
    yield put(alertActions.addAlert('Viewing Game!', 'success'));
    yield put(push(`/game/${gameID}`));
  } catch (err) {
    yield put({
      type: game.reducer.VIEW_GAME_FAILED,
      payload: err.response.data,
    });
    yield put(alertActions.addAlert(err.response.data, 'error'));
  }
}
