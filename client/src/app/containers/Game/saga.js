import { all, put, select, takeLatest, call } from 'redux-saga/effects';
import axios from 'axios';
import { selectCurrentUser } from '../../../auth/selectors';
import { selectCurrentGameID, selectCurrentAction } from './selectors';
import { actions as gameActions } from './slice';
import { actions as alertActions } from '../Alert/slice';
import { alertTypes } from '../Alert/types';

export function* handleExitGame({ payload: { history } }) {
  yield call(history.push, '/home');
}

export function* handleGoToGame({ payload: { id, history } }) {
  if (!id) {
    yield put(
      alertActions.addAlert('No id in handleGoToGame', alertTypes.error),
    );
    return;
  }

  const currentUser = yield select(selectCurrentUser);

  // select the id being used to login from the state
  try {
    const res = yield axios.get(`/game/${id}?player=${currentUser.id}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
    yield call(history.push, '/game');
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `something bad happened... ${err}`,
        alertTypes.error,
      ),
    );
  }
}

export function* handleRefreshCurrentGame({ payload: { id, history } }) {
  const currentUser = yield select(selectCurrentUser);

  try {
    const res = yield axios.get(`/game/${id}?player=${currentUser.id}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* handleDeal({ payload: { history } }) {
  const currentUser = yield select(selectCurrentUser);
  const gameID = yield select(selectCurrentGameID);
  const currentAction = yield select(selectCurrentAction);
  console.log(`num shuffles: ${currentAction.numShuffles}`);

  try {
    // the return of the post is just 'action handled'
    // it may be wise to make it return the new game state?
    yield axios.post('/action', {
      pID: currentUser.id,
      gID: gameID,
      o: 0, // magic number means "deals cards"
      a: { ns: currentAction.numShuffles },
    });
    // TODO wait a moment and re-fetch?
    yield put(gameActions.refreshGame({ id: gameID }));
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `handling deal broke ${err.response ? err.response.data : err}`,
        alertTypes.error,
      ),
    );
  }
}

export function* gameSaga() {
  yield all([
    takeLatest(gameActions.goToGame.type, handleGoToGame),
    takeLatest(gameActions.exitGame.type, handleExitGame),
    takeLatest(gameActions.refreshGame.type, handleRefreshCurrentGame),
    takeLatest(gameActions.dealCards.type, handleDeal),
  ]);
}
