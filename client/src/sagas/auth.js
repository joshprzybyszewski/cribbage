import { delay, put, takeLatest } from 'redux-saga/effects';
import { LOGIN, LOGIN_ASYNC } from './types';

// Worker saga for login
export function* loginAsync({ payload }) {
  yield delay(1000);
  yield put({ type: LOGIN, payload });
}

// Watcher saga for login
export function* watchLoginAsync() {
  yield takeLatest(LOGIN_ASYNC, loginAsync);
}
