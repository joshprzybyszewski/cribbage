import { all, put, takeLatest } from 'redux-saga/effects';
import axios from 'axios';
import { actions as homeActions } from './slice';
import { actions as alertActions } from '../Alert/slice';

export function* handleRefreshActiveGames({ payload: { id } }) {
  if (!id) {
    yield put(alertActions.addAlert('undefined player ID', 'warning'));
    return;
  }

  try {
    const res = yield axios.get(`/games/active?playerID=${id}`);
    const { player, activeGames } = res.data;
    yield put(homeActions.gotActiveGames({ player, activeGames }));
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, 'error'));
  }
}

export function* homeSaga() {
  yield all([
    takeLatest(homeActions.refreshActiveGames.type, handleRefreshActiveGames),
  ]);
}
