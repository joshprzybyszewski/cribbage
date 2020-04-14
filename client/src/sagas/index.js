import { all } from 'redux-saga/effects';
import { watchLoginAsync } from './auth';

export default function* rootSaga() {
  yield all([watchLoginAsync()]);
}
