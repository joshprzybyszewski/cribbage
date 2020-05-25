import { delay, put, takeLatest } from 'redux-saga/effects';
import { ADD_ALERT_TRIGGER, ADD_ALERT } from './types';
import { uuid } from 'uuidv4';

export const addAlertAction = (msg, type) => ({
  type: ADD_ALERT_TRIGGER,
  payload: { msg, type },
});

export function* handleAddAlert({ payload }) {
  const id = uuid();
  yield put({
    type: ADD_ALERT,
    payload: {
      id,
      msg: payload.msg,
      type: payload.type,
    },
  });
  yield delay(5000);
  yield put({
    type: REMOVE_ALERT,
    payload: id,
  });
}
export function* watchAddAlert() {
  yield takeLatest(ADD_ALERT_TRIGGER, handleAddAlert);
}
