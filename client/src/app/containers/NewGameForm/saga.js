import { actions as alertActions } from 'app/containers/Alert/slice';
import { alertTypes } from 'app/containers/Alert/types';
import { actions as newGameActions } from 'app/containers/NewGameForm/slice';
import { selectCurrentUser } from 'auth/selectors';
import axios from 'axios';
import { all, call, put, select, takeLatest } from 'redux-saga/effects';

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
    yield call(history.push, `/game/${id}`);
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* createGameSaga() {
  yield all([takeLatest(newGameActions.createGame.type, handleCreateGame)]);
}
