import {
    selectHandCards,
} from 'app/containers/Suggestions/selectors';
import { actions as sugActions } from 'app/containers/Suggestions/slice';
import axios from 'axios';
import { all, put, select, takeLatest, call } from 'redux-saga/effects';


export function* handleSuggestHandRequest() {
    const currentHand = yield select(selectHandCards);

    try {
        const res = yield axios.get(`/suggest/hand?dealt=${currentHand}`);
        yield put(sugActions.setSuggestionResult({ data: res.data }));
    } catch (err) {
        throw err
    }
}

export function* suggestionsSaga() {
    yield all([
        takeLatest(sugActions.getHandSuggestion.type, handleSuggestHandRequest),
    ]);
}