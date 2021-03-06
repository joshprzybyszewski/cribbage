import { actions } from 'app/containers/Alert/slice';
import { delay, put, takeEvery } from 'redux-saga/effects';

export function* handleAlert({ payload: { id } }) {
    yield delay(5000);
    yield put(actions.removeAlert(id));
}
export function* alertSaga() {
    yield takeEvery(actions.addAlert.type, handleAlert);
}
