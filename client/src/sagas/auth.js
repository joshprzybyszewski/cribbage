import { put, takeLatest, call } from 'redux-saga/effects';
import {
  LOGIN_ASYNC,
  LOGIN_FAIL,
  LOGIN_SUCCESS,
  REGISTER_ASYNC,
  REGISTER_FAIL,
  REGISTER_SUCCESS,
} from './types';
import axios from 'axios';
import { setAlert } from './alert';

export const login = username => ({
  type: LOGIN_ASYNC,
  payload: username,
});

export function* loginAsync({ payload }) {
  try {
    console.log(`endpoint: /player/${payload}`);
    const res = yield call(axios.get, `/player/${payload}`);
    yield put({ type: LOGIN_SUCCESS, payload: res.data });
    yield put(setAlert('Successfully logged in!', 'success'));
  } catch (err) {
    yield put(setAlert(err.response.data, 'error'));
    yield put({ type: LOGIN_FAIL, payload: err.response.data });
  }
}

export const register = (username, displayName) => ({
  type: REGISTER_ASYNC,
  payload: { username, displayName },
});

export function* registerAsync({ payload: { username, displayName } }) {
  try {
    const res = yield call(
      axios.post,
      `/create/player/${username}/${displayName}`
    );
    yield put({ type: REGISTER_SUCCESS, payload: res.data });
    yield put(setAlert('Successfully registered!', 'success'));
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
