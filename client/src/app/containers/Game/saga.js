import { actions as alertActions } from 'app/containers/Alert/slice';
import { alertTypes } from 'app/containers/Alert/types';
import { phase } from 'app/containers/Game/constants';
import { newPlayerAction } from 'app/containers/Game/convert';
import { actions as gameActions } from 'app/containers/Game/slice';
import { selectCurrentUser } from 'auth/selectors';
import axios from 'axios';
import { all, put, select, takeLatest, call } from 'redux-saga/effects';

const cardToGolangCard = c => {
  const magicMap = {
    Spades: 0,
    Clubs: 1,
    Diamonds: 2,
    Hearts: 3,
  };
  return {
    s: magicMap[c.suit],
    v: c.value,
  };
};

export function* handleRequestGame({ payload: { userID, gameID } }) {
  try {
    const res = yield axios.get(`/game/${gameID}?player=${userID}`);
    yield put(gameActions.requestGameSuccess(res.data));
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `something bad happened... ${err}`,
        alertTypes.error,
      ),
    );
    yield put(gameActions.requestGameFailure());
  }
}

export function* handleExitGame({ payload: { history } }) {
  yield call(history.push, '/home');
}

export function* handleRefreshCurrentGame({ payload: { id } }) {
  const currentUser = yield select(selectCurrentUser);

  try {
    const res = yield axios.get(`/game/${id}?player=${currentUser.id}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
  } catch (err) {
    yield put(alertActions.addAlert(err.response.data, alertTypes.error));
  }
}

// postAction returns the next redux action to dispatch so each function* can `put` it
const postAction = async playerAction => {
  try {
    await axios.post('/action', playerAction);
  } catch (err) {
    return alertActions.addAlert(
      `handling action broke ${err.response ? err.response.data : err}`,
      alertTypes.error,
    );
  }
  return gameActions.refreshGame(playerAction.gID);
};
export function* handleBuildCrib({ payload: { userID, gameID, cards } }) {
  const playerAction = newPlayerAction(userID, gameID, phase.crib, {
    cs: cards.map(c => cardToGolangCard(c)),
  });
  const next = yield postAction(playerAction);
  yield put(gameActions.clearSelectedCards());
  yield put(next);
}

export function* handlePeg({ payload: { userID, gameID, card } }) {
  const playerAction = newPlayerAction(userID, gameID, phase.peg, {
    c: cardToGolangCard(card),
  });
  const next = yield postAction(playerAction);
  yield put(gameActions.clearSelectedCards());
  yield put(next);
}

export function* handleSayGo({ payload: { userID, gameID } }) {
  const playerAction = newPlayerAction(userID, gameID, phase.peg, { sg: true });
  const next = yield postAction(playerAction);
  yield put(next);
}

export function* handleDeal({ payload: { userID, gameID, numShuffles } }) {
  const playerAction = newPlayerAction(userID, gameID, phase.deal, {
    ns: numShuffles,
  });
  const next = yield postAction(playerAction);
  yield put(next);
}

export function* handleCutDeck({ payload: { userID, gameID, cutPct } }) {
  const playerAction = newPlayerAction(userID, gameID, phase.cut, {
    p: cutPct,
  });
  const next = yield postAction(playerAction);
  yield put(next);
}

export function* handleCountHand({
  payload: { userID, gameID, points, isCrib },
}) {
  const playerAction = newPlayerAction(
    userID,
    gameID,
    isCrib ? phase.countCrib : phase.countHand,
    {
      pts: points,
    },
  );
  const next = yield postAction(playerAction);
  yield put(next);
}

export function* gameSaga() {
  yield all([
    takeLatest(gameActions.requestGame.type, handleRequestGame),
    takeLatest(gameActions.exitGame.type, handleExitGame),
    takeLatest(gameActions.refreshGame.type, handleRefreshCurrentGame),
    takeLatest(gameActions.dealCards.type, handleDeal),
    takeLatest(gameActions.buildCrib.type, handleBuildCrib),
    takeLatest(gameActions.cutDeck.type, handleCutDeck),
    takeLatest(gameActions.pegCard.type, handlePeg),
    takeLatest(gameActions.sayGo.type, handleSayGo),
    takeLatest(gameActions.countHand.type, handleCountHand),
  ]);
}
