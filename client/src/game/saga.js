import { all, put, select, takeLatest, call } from 'redux-saga/effects';
import axios from 'axios';
import { selectCurrentUser } from '../auth/selectors';
import { selectCurrentGameID } from './selectors';
import { actions as gameActions } from './slice';
import { actions as alertActions } from '../app/containers/Alert/slice';

export function* handleExitGame({ payload: { history } }) {
  yield call(history.push, '/home');
}

export function* handleGoToGame({ payload: { id, history } }) {
  if (!id) {
    yield put(
      alertActions.addAlert(
        'No id in handleGoToGame',
        'could not get game to go to',
      ),
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
        'error could not get game',
      ),
    );
  }
}

export function* handleRefreshCurrentGame({ payload: { history } }) {
  const currentGameID = yield select(selectCurrentGameID);
  if (!currentGameID) {
    yield put(
      alertActions.addAlert(
        'No currentGameID',
        'could not refresh current game',
      ),
    );
    return;
  }

  const currentUser = yield select(selectCurrentUser);

  try {
    const res = yield axios.get(
      `/game/${currentGameID}?player=${currentUser.id}`,
    );
    yield put(gameActions.gameRetrieved({ data: res.data }));
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
