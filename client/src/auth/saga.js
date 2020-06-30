import { push } from 'connected-react-router';
import { all, put, takeLatest } from 'redux-saga/effects';
import axios from 'axios';

import { actions as authActions } from './slice';
import { actions as alertActions } from '../app/containers/Alert/slice';

export function* handleLogout() {
  yield put(push('/'));
}

export function* handleLogin({ payload }) {
  try {
    const res = yield axios.get(`/player/${payload}`);
    yield put(
      authActions.loginSuccess({
        id: res.data.player.id,
        name: res.data.player.name,
      }),
    );
    yield put(
      alertActions.addAlert({ msg: 'Login successful!', type: 'success' }),
    );
    yield put(push('/home'));
  } catch (err) {
    yield put(authActions.loginFailed(err.response.data));
    yield put(alertActions.addAlert({ msg: err.response.data, type: 'error' }));
  }
}

export function* handleRegister({ payload: { id, name } }) {
  try {
    const res = yield axios.post('/create/player', { player: { id, name } });
    yield put(
      authActions.registerSuccess({
        id: res.data.player.id,
        name: res.data.player.name,
      }),
    );
    yield put(
      alertActions.addAlert({
        msg: 'Registration successful!',
        type: 'success',
      }),
    );
    yield put(push('/home'));
  } catch (err) {
    yield put(authActions.registerFailed(err.response.data));
    yield put(alertActions.addAlert({ msg: err.response.data, type: 'error' }));
  }
}

export function* authSaga() {
  yield all([
    takeLatest(authActions.login.type, handleLogin),
    takeLatest(authActions.register.type, handleRegister),
    takeLatest(authActions.logout.type, handleLogout),
  ]);
}
