import { push } from 'connected-react-router';
import { put, takeLatest } from 'redux-saga/effects';
import {
  LOGIN,
  LOGIN_ASYNC,
  LOGIN_FAILED,
  LOGOUT,
  LOGOUT_TRIGGER,
  REGISTER,
  REGISTER_ASYNC,
  REGISTER_FAILED,
} from './types';
import axios from 'axios';
import { addAlertAction } from './alert';

export const logoutAction = () => ({
  type: LOGOUT_TRIGGER,
});

export const loginAction = id => ({
  type: LOGIN_ASYNC,
  payload: id,
});

export function* logout() {
  yield put({ type: LOGOUT });
}

export function* loginAsync({ payload }) {
  try {
    const res = yield axios.get(`/player/${payload}`);
    yield put({
      type: LOGIN,
      payload: { id: res.data.id, name: res.data.name },
    });
    yield put(addAlertAction('Login successful!', 'success'));
    yield put(push('/home'));
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
    yield put(push('/home'));
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
export function* watchLogout() {
  yield takeLatest(LOGOUT_TRIGGER, logout);
}

export function* watchRegisterAsync() {
  yield takeLatest(REGISTER_ASYNC, registerAsync);
}
