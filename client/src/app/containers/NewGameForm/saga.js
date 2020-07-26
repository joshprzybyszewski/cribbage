import axios from 'axios';
import { all, put, select, takeLatest, call } from 'redux-saga/effects';

import { actions as alertActions } from '../Alert/slice';
import { actions as gameActions } from '../Game/slice';
import { alertTypes } from '../Alert/types';
import { selectCurrentUser } from '../../../auth/selectors';
import { actions as newGameActions } from './slice';

export function* handleCreateGame({
  payload: { opp1ID, opp2ID, teammateID, history },
}) {
  const currentUser = yield select(selectCurrentUser);
  try {
    let playerIDs = [currentUser.id, opp1ID];
    if (teammateID) {
      playerIDs.push(teammateID);
    }
    if (opp2ID) {
      playerIDs.push(opp2ID);
    }
    const res = yield axios.post(`/create/game`, {
      playerIDs: playerIDs,
    });
    const id = res.data.id;
    yield put(gameActions.goToGame(id, history));
    // yield call(history.push, '/home');
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* createGameSaga() {
  yield all([takeLatest(newGameActions.createGame.type, handleCreateGame)]);
}
