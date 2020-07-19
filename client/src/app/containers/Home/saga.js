import axios from 'axios';
import { all, put, takeLatest } from 'redux-saga/effects';

import { actions as alertActions } from '../Alert/slice';
import { alertTypes } from '../Alert/types';
import { actions as homeActions } from './slice';

export function* handleRefreshActiveGames({ payload: { id } }) {
  if (!id) {
    yield put(alertActions.addAlert('undefined player ID', alertTypes.warning));
    return;
  }

  try {
    const res = yield axios.get(`/games/active?playerID=${id}`);
    const { player, activeGames } = res.data;
    yield put(homeActions.gotActiveGames({ player, activeGames }));
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* homeSaga() {
  yield all([
    takeLatest(homeActions.refreshActiveGames.type, handleRefreshActiveGames),
  ]);
}
