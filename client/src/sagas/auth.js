import { put, takeLatest } from 'redux-saga/effects';
import {
  LOGIN,
  LOGIN_ASYNC,
  LOGIN_FAILED,
  REGISTER,
  REGISTER_ASYNC,
  REGISTER_FAILED,
} from './types';
import axios from 'axios';
import { addAlertAction } from './alert';

export const loginAction = id => ({
  type: LOGIN_ASYNC,
  payload: id,
});

export function* loginAsync({ payload }) {
  try {
    const res = yield axios.get(`/player/${payload}`);
    yield put({
      type: LOGIN,
      payload: { id: res.data.id, name: res.data.name },
    });
    yield put(addAlertAction('Login successful!', 'success'));
  } catch (err) {
    yield put({
      type: LOGIN_FAILED,
      payload: err.response.data,
    });
    yield put(addAlertAction(err.response.data, 'error'));
  }
}

export const registerAction = (id, name) => ({
  type: REGISTER_ASYNC,
  payload: { id, name },
});

export function* registerAsync({ payload: { id, name } }) {
  try {
    const res = yield axios.post('/create/player', { id, name });
    yield put({
      type: REGISTER,
      payload: { id: res.data.id, name: res.data.name },
    });
    yield put(addAlertAction('Registration successful!', 'success'));
  } catch (err) {
    yield put({
      type: REGISTER_FAILED,
      payload: err.response.data,
    });
    yield put(addAlertAction(err.response.data, 'error'));
  }
}

export function* watchLoginAsync() {
  yield takeLatest(LOGIN_ASYNC, loginAsync);
}

export function* watchRegisterAsync() {
  yield takeLatest(REGISTER_ASYNC, registerAsync);
}
