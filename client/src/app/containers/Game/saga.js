import { all, put, select, takeLatest, call } from 'redux-saga/effects';
import axios from 'axios';
import { selectCurrentUser } from '../../../auth/selectors';
import { selectCurrentGameID, selectCurrentAction } from './selectors';
import { actions as gameActions } from './slice';
import { actions as alertActions } from '../Alert/slice';
import { alertTypes } from '../Alert/types';

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

  // select the id being used to login from the state
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

export function* handleRefreshCurrentGame({ payload: { id, history } }) {
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

const getPlayerActionJSON = (myID, gID, phase, currentAction) => {
  let overcomesMap = {
    deal: 0,
    crib: 1,
    cut: 2,
    peg: 3,
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
      console.log('is pegging?');
      action =
        !currentAction.selectedCards || currentAction.selectedCards.length !== 1
          ? {
              sg: true,
            }
          : {
              c: cardToGolangCard(currentAction.selectedCards[0]),
            };
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
      getPlayerActionJSON(currentUser.id, id, phase, currentAction),
    );
    // TODO wait a moment and re-fetch?
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

export function* handleCutDeck() {
  yield handleGenericAction('cut');
}

export function* handPegCard() {
  yield handleGenericAction('peg');
}

export function* gameSaga() {
  yield all([
    takeLatest(gameActions.goToGame.type, handleGoToGame),
    takeLatest(gameActions.exitGame.type, handleExitGame),
    takeLatest(gameActions.refreshGame.type, handleRefreshCurrentGame),
    takeLatest(gameActions.dealCards.type, handleDeal),
    takeLatest(gameActions.buildCrib.type, handleBuildCrib),
    takeLatest(gameActions.cutDeck.type, handleCutDeck),
    takeLatest(gameActions.pegCard.type, handPegCard),
  ]);
}
