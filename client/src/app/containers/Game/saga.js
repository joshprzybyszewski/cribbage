import { actions as alertActions } from 'app/containers/Alert/slice';
import { alertTypes } from 'app/containers/Alert/types';
import { newPlayerAction } from 'app/containers/Game/convert';
import {
  selectCurrentGameID,
  selectCurrentAction,
} from 'app/containers/Game/selectors';
import { actions as gameActions } from 'app/containers/Game/slice';
import { selectCurrentUser } from 'auth/selectors';
import axios from 'axios';
import { all, put, select, takeLatest, call } from 'redux-saga/effects';

export function* handleExitGame({ payload: { history } }) {
  yield call(history.push, '/home');
}

export function* handleGoToGame({ payload: { id, history } }) {
  if (!id) {
    yield put(
      alertActions.addAlert('No id in handleGoToGame', alertTypes.error),
    );
    return;
  }

  const currentUser = yield select(selectCurrentUser);

  try {
    const res = yield axios.get(`/game/${id}?player=${currentUser.id}`);
    yield put(gameActions.gameRetrieved({ data: res.data }));
    yield call(history.push, '/game');
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `something bad happened... ${err}`,
        alertTypes.error,
      ),
    );
  }
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

// getPlayerAction returns the JSON struct which the server knows
// how to interpret
const getPlayerAction = (myID, gID, phase, currentAction) => {
  const overcomesMap = {
    deal: 0,
    crib: 1,
    cut: 2,
    peg: 3,
    counthand: 4,
    countcrib: 5,
  };
  let action = {};
  switch (phase) {
    case 'deal':
      action = { ns: currentAction.numShuffles };
      break;
    case 'crib':
      action = {
        cs: currentAction.selectedCards
          ? currentAction.selectedCards.map(cardToGolangCard)
          : [],
      };
      break;
    case 'cut':
      action = { p: currentAction.percCut };
      break;
    case 'peg':
      action =
        !currentAction.selectedCards || currentAction.selectedCards.length !== 1
          ? {
              sg: true,
            }
          : {
              c: cardToGolangCard(currentAction.selectedCards[0]),
            };
      break;
    case 'counthand':
      action = { pts: currentAction.points };
      break;
    case 'countcrib':
      action = { pts: currentAction.points };
      break;
    default:
      action = { badstate: true };
      break;
  }
  return {
    pID: myID,
    gID: gID,
    o: overcomesMap[phase],
    a: action,
  };
};

function* handleGenericAction(phase) {
  const currentUser = yield select(selectCurrentUser);
  const id = yield select(selectCurrentGameID);
  const currentAction = yield select(selectCurrentAction);

  try {
    // the return of the post is just 'action handled'
    // it may be wise to make it return the new game state?
    yield axios.post(
      '/action',
      getPlayerAction(currentUser.id, id, phase, currentAction),
    );
    yield put(gameActions.refreshGame(id));
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `handling action broke ${err.response ? err.response.data : err}`,
        alertTypes.error,
      ),
    );
  }
}

export function* handleDeal() {
  yield handleGenericAction('deal');
}

export function* handleBuildCrib() {
  yield handleGenericAction('crib');
}

export function* handleCutDeck({ payload: cutPct }) {
  const currentUser = yield select(selectCurrentUser);
  const gameID = yield select(selectCurrentGameID);
  const playerAction = newPlayerAction(currentUser.id, gameID, 'cut', {
    p: cutPct,
  });
  try {
    yield axios.post('/action', playerAction);
    yield put(gameActions.refreshGame(gameID));
  } catch (err) {
    yield put(
      alertActions.addAlert(
        `handling action broke ${err.response ? err.response.data : err}`,
        alertTypes.error,
      ),
    );
  }
}

export function* handlePeg() {
  yield handleGenericAction('peg');
}

export function* handleCountHand() {
  yield handleGenericAction('counthand');
}

export function* handleCountCrib() {
  yield handleGenericAction('countcrib');
}

export function* gameSaga() {
  yield all([
    takeLatest(gameActions.goToGame.type, handleGoToGame),
    takeLatest(gameActions.exitGame.type, handleExitGame),
    takeLatest(gameActions.refreshGame.type, handleRefreshCurrentGame),
    takeLatest(gameActions.dealCards.type, handleDeal),
    takeLatest(gameActions.buildCrib.type, handleBuildCrib),
    takeLatest(gameActions.cutDeck.type, handleCutDeck),
    takeLatest(gameActions.pegCard.type, handlePeg),
    takeLatest(gameActions.countHand.type, handleCountHand),
    takeLatest(gameActions.countCrib.type, handleCountCrib),
  ]);
}
