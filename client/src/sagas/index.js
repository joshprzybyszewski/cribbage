import { all } from 'redux-saga/effects';
import { watchLoginAsync, watchRegisterAsync } from './auth';
import { watchAddAlert } from './alert';

export default function* rootSaga() {
  yield all([watchLoginAsync(), watchRegisterAsync(), watchAddAlert()]);
}
