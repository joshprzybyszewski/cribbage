import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentGameID: '',
  currentGame: {},
  selectedCards: [],
  isLoading: true,
};

const gameSlice = createSlice({
  name: 'game',
  initialState,
  reducers: {
    requestGame(state, action) {
      state.currentGameID = action.payload.gameID;
      state.isLoading = true;
    },
    requestGameSuccess(state, action) {
      state.currentGame = action.payload;
      state.isLoading = false;
    },
    requestGameFailure(state) {
      state.currentGame = {};
      state.isLoading = false;
    },
    gameRetrieved(state, action) {
      state.isLoading = false;
      state.currentGame = action.payload.data;
      state.currentAction = initialState.currentAction;
      switch (state.currentGame.phase) {
        case `Deal`:
          // TODO leave numShuffles
          break;
        default:
          // TODO here too
          break;
      }
    },
    exitGame: {
      reducer: state => {
        state.isLoading = false;
        state.currentGameID = '';
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    selectCard(state, action) {
      const { payload: card } = action;
      if (!state.selectedCards.map(c => c.name).includes(card.name)) {
        state.selectedCards.push(card);
      }
    },
    unselectCard(state, action) {
      const { payload: card } = action;
      state.selectedCards = state.selectedCards.filter(
        c => c.name !== card.name,
      );
    },
    clearSelectedCards(state) {
      state.selectedCards = [];
    },
    dealCards() {},
    buildCrib() {},
    cutDeck() {},
    pegCard() {},
    sayGo() {},
    countHand() {},
  },
});

export const { actions, reducer, name: sliceKey } = gameSlice;
