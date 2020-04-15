import { delay, put, takeLatest, call } from 'redux-saga/effects';
import {
  LOGIN,
  LOGIN_ASYNC,
  REGISTER_FAIL,
  REGISTER_SUCCESS,
  REGISTER_ASYNC,
} from './types';
import axios from 'axios';
import { setAlert } from './alert';

export function* loginAsync({ payload }) {
  yield delay(1000);
  yield put({ type: LOGIN, payload });
}

export function* registerAsync({ payload }) {
  const { username, displayName } = payload;
  try {
    const res = yield call(
      axios.post,
      `/create/player/${username}/${displayName}`
    );
    yield put({ type: REGISTER_SUCCESS, payload: res.data });
  } catch (err) {
    yield put(setAlert(err.response.data, 'error'));
    yield put({ type: REGISTER_FAIL, payload: err.response.data });
  }
}

export function* watchLoginAsync() {
  yield takeLatest(LOGIN_ASYNC, loginAsync);
}

export function* watchRegisterAsync() {
  yield takeLatest(REGISTER_ASYNC, registerAsync);
}
