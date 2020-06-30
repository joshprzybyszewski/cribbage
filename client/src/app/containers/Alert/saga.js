import { delay, put, takeEvery } from 'redux-saga/effects';
import { uuid } from 'uuidv4';
import { actions } from './slice';

export function* addAlert({ payload: { msg, type } }) {
  const id = uuid();
  yield put(actions.addAlert({ id, msg, type }));
  yield delay(5000);
  yield put(actions.removeAlert(id));
}
export function* watchAddAlert() {
  yield takeEvery(actions.requestAlert(), addAlert);
}
