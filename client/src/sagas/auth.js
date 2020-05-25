import { delay, put, takeLatest } from 'redux-saga/effects';
import {
  LOGIN,
  LOGIN_ASYNC,
  REGISTER,
  REGISTER_ASYNC,
  REGISTER_FAILED,
} from './types';
import axios from 'axios';

export function* loginAsync({ payload }) {
  yield delay(1000);
  yield put({ type: LOGIN, payload });
}

export function* registerAsync({ payload: { id, name } }) {
  try {
    const res = yield axios.post('/create/player', { id, name });
    yield put({
      type: REGISTER,
      payload: { id: res.data.id, name: res.data.name },
    });
  } catch (err) {
    yield put({
      type: REGISTER_FAILED,
      payload: err.response.data,
    });
  }
}

export function* watchLoginAsync() {
  yield takeLatest(LOGIN_ASYNC, loginAsync);
}

export function* watchRegisterAsync() {
  yield takeLatest(REGISTER_ASYNC, registerAsync);
}
