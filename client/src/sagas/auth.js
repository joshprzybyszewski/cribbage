import { push } from 'connected-react-router';
import { put, takeLatest } from 'redux-saga/effects';
import axios from 'axios';

import { auth } from './types';
import { alertActions } from './actions';

// logout
export function* watchLogout() {
  yield takeLatest(auth.LOGOUT, logout);
}

export function* logout() {
  yield put({ type: auth.reducer.LOGOUT });
  yield put(push('/'));
}

// login
export function* watchLoginAsync() {
  yield takeLatest(auth.LOGIN, loginAsync);
}

export function* loginAsync({ payload }) {
  try {
    let pID = payload;
    // const res = yield axios.get(`/player/${pID}`);
    yield put(alertActions.addAlert('Login successful!', 'success'));
    const allgamesRes = yield axios.get(`/player/${pID}/allgames`);
    // console.log(allgamesRes);
    yield put({
      type: auth.reducer.LOGIN,
      payload: {
        // id: res.data.player.id,
        // name: res.data.player.name,
        id: allgamesRes.data.player.id,
        name: allgamesRes.data.player.name,
        allgames: allgamesRes.data.allgames,
      },
    });
    yield put(alertActions.addAlert('Allgames gotten successfully!', 'success'));
    yield put(push('/home'));
  } catch (err) {
    yield put({
      type: auth.reducer.LOGIN_FAILED,
      payload: err.response.data,
    });
    yield put(alertActions.addAlert(err.response.data, 'error'));
  }
}

// register
export function* watchRegisterAsync() {
  yield takeLatest(auth.REGISTER, registerAsync);
}

export function* registerAsync({ payload: { id, name } }) {
  try {
    const res = yield axios.post('/create/player', { player: { id, name } });
    yield put({
      type: auth.reducer.REGISTER,
      payload: { id: res.data.player.id, name: res.data.player.name },
    });
    yield put(alertActions.addAlert('Registration successful!', 'success'));
    yield put(push('/home'));
  } catch (err) {
    yield put({
      type: auth.reducer.REGISTER_FAILED,
      payload: err.response.data,
    });
    yield put(alertActions.addAlert(err.response.data, 'error'));
  }
}
