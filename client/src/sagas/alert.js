import { delay, put, takeEvery } from 'redux-saga/effects';
import { alert } from './types';
import { uuid } from 'uuidv4';

export function* handleAddAlert({ payload }) {
  const id = uuid();
  yield put({
    type: alert.reducer.ADD_ALERT,
    payload: {
      id,
      msg: payload.msg,
      type: payload.type,
    },
  });
  yield delay(5000);
  yield put({
    type: alert.reducer.REMOVE_ALERT,
    payload: id,
  });
}
export function* watchAddAlert() {
  yield takeEvery(alert.ADD_ALERT, handleAddAlert);
}
