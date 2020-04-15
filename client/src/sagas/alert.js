import { delay, put, takeEvery } from 'redux-saga/effects';
import { v4 as uuidv4 } from 'uuid';
import { REMOVE_ALERT, SET_ALERT, SET_ALERT_WORKER } from './types';

export const setAlert = (message, alertType) => {
  const id = uuidv4();
  return {
    type: SET_ALERT,
    payload: { id, message, type: alertType },
  };
};

export function* setAlertHandler({ payload }) {
  yield put({ type: SET_ALERT_WORKER, payload });
  yield delay(5000);
  yield put({ type: REMOVE_ALERT, payload: payload.id });
}

export function* watchSetAlert() {
  yield takeEvery(SET_ALERT, setAlertHandler);
}
