import { all } from 'redux-saga/effects';
import { watchLoginAsync, watchLogout, watchRegisterAsync } from './auth';
import { watchAddAlert } from './alert';

export default function* rootSaga() {
  yield all([
    watchLoginAsync(),
    watchLogout(),
    watchRegisterAsync(),
    watchAddAlert(),
  ]);
}
