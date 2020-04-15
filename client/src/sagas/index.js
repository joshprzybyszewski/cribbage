import { all } from 'redux-saga/effects';
import { watchSetAlert } from './alert';
import { watchLoginAsync, watchRegisterAsync } from './auth';

export default function* rootSaga() {
  yield all([watchLoginAsync(), watchRegisterAsync(), watchSetAlert()]);
}
