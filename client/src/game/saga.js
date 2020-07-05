import { all, put, select, takeLatest, call } from 'redux-saga/effects';
import axios from 'axios';
import { selectCurrentGameID } from './selectors';
import { actions as gameActions } from './slice';
import { actions as alertActions } from '../app/containers/Alert/slice';

export function* handleExitGame({ payload: { history } }) {
  yield call(history.push, '/home');
}

export function* handleGoToGame({ payload: { id, history } }) {
  if (!id) {
    yield put(alertActions.addAlert('No id in handleGoToGame', 'could not get game to go to'));
    return;
  }

  // select the id being used to login from the state
  try {
    const res = yield axios.get(`/game/${id}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
    yield call(history.push, '/game');
    yield put(alertActions.addAlert('Game Got!', 'success'));
} catch (err) {
    yield put(alertActions.addAlert(err.response.data, 'error could not get game'));
  }
}

export function* handleRefreshCurrentGame({ payload: { history } }) {
    
  const currentGameID = yield select(selectCurrentGameID);
  if (!currentGameID) {
    yield put(alertActions.addAlert('No currentGameID', 'could not refresh current game'));
    return;
  }

  try {
    const res = yield axios.get(`/game/${currentGameID}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
    yield put(alertActions.addAlert('Game Refreshed!', 'success'));
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, 'error'));
  }
}

export function* gameSaga() {
  yield all([
    takeLatest(gameActions.goToGame.type, handleGoToGame),
    takeLatest(gameActions.exitGame.type, handleExitGame),
    takeLatest(gameActions.refreshGame.type, handleRefreshCurrentGame),
  ]);
}
