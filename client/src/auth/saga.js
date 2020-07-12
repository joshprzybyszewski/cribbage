import { all, put, select, takeLatest, call } from 'redux-saga/effects';
import axios from 'axios';
import { selectCurrentUser } from './selectors';
import { actions as authActions } from './slice';
import { actions as alertActions } from '../app/containers/Alert/slice';
import { alertTypes } from '../app/containers/Alert/types';
import { actions as homeActions } from '../app/containers/Home/slice';

export function* handleLogout({ payload: { history } }) {
  yield call(history.push, '/');
}

export function* handleLogin({ payload: { history } }) {
  // select the id being used to login from the state
  const currentUser = yield select(selectCurrentUser);
  try {
    const res = yield axios.get(`/player/${currentUser.id}`);
    const { id, name } = res.data.player;
    yield put(authActions.loginSuccess({ id, name }));
    yield put(homeActions.refreshActiveGames({ id: currentUser.id }));
    yield call(history.push, '/home');
  } catch (err) {
    yield put(authActions.loginFailed(err.response.data));
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* handleRegister({ payload: { history } }) {
  // select the id and name being used to register from the state
  const currentUser = yield select(selectCurrentUser);
  try {
    const res = yield axios.post('/create/player', { player: currentUser });
    const { id, name } = res.data.player;
    yield put(authActions.registerSuccess({ id, name }));
    yield put(
      alertActions.addAlert('Registration successful!', alertTypes.success),
    );
    yield call(history.push, '/home');
  } catch (err) {
    yield put(authActions.registerFailed(err.response.data));
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

export function* authSaga() {
  yield all([
    takeLatest(authActions.login.type, handleLogin),
    takeLatest(authActions.register.type, handleRegister),
    takeLatest(authActions.logout.type, handleLogout),
  ]);
}
