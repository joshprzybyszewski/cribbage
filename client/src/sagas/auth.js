import { delay, put, takeLatest } from 'redux-saga/effects';
import { LOGIN, LOGIN_ASYNC, REGISTER, REGISTER_ASYNC } from './types';

export function* loginAsync({ payload }) {
  yield delay(1000);
  yield put({ type: LOGIN, payload });
}

export function* registerAsync({ payload }) {
  yield delay(1000);
  yield put({ type: REGISTER, payload });
}

export function* watchLoginAsync() {
  yield takeLatest(LOGIN_ASYNC, loginAsync);
}

export function* watchRegisterAsync() {
  yield takeLatest(REGISTER_ASYNC, registerAsync);
}
