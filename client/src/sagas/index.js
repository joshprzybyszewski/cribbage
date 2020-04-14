import { all } from 'redux-saga/effects';
import { watchLoginAsync, watchRegisterAsync } from './auth';

export default function* rootSaga() {
  yield all([watchLoginAsync(), watchRegisterAsync()]);
}
