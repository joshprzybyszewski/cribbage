import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentGame: {},
  selectedCards: [],
  isLoading: true,
};

const gameSlice = createSlice({
  name: 'game',
  initialState,
  reducers: {
    requestGame(state) {
      state.isLoading = true;
    },
    requestGameSuccess(state, action) {
      state.isLoading = false;
      state.currentGame = action.payload;
    },
    requestGameFailure(state) {
      state.isLoading = false;
      state.currentGame = {};
    },
    // eslint-disable-next-line no-unused-vars
    exitGame(state) {
      state = initialState;
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
