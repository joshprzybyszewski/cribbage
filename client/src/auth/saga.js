import { push } from 'connected-react-router';
import axios from 'axios';
import { all, put, select, takeLatest } from 'redux-saga/effects';

import { selectCurrentUser } from './selectors';
import { actions as authActions } from './slice';
import { actions as alertActions } from '../app/containers/Alert/slice';

export function* handleLogout() {
  yield put(push('/'));
}

export function* handleLogin() {
  // select the id being used to login from the state
  const currentUser = yield select(selectCurrentUser);
  try {
    const res = yield axios.get(`/player/${currentUser.id}`);
    const { id, name } = res.data.player;
    yield put(authActions.loginSuccess({ id, name }));
    yield put(
      alertActions.addAlert({ msg: 'Login successful!', type: 'success' }),
    );
    yield put(push('/home'));
  } catch (err) {
    yield put(authActions.loginFailed(err.response.data));
    yield put(alertActions.addAlert({ msg: err.response.data, type: 'error' }));
  }
}

export function* handleRegister() {
  // select the id and name being used to register from the state
  const currentUser = yield select(selectCurrentUser);
  try {
    const res = yield axios.post('/create/player', { player: currentUser });
    const { id, name } = res.data.player;
    yield put(authActions.registerSuccess({ id, name }));
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
