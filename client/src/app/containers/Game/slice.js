import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentGameID: '',
  currentGame: {},
  currentAction: {},
  loading: true,
};

const gameSlice = createSlice({
  name: 'game',
  initialState,
  reducers: {
    goToGame: {
      reducer: (state, action) => {
        state.loading = true;
        state.currentGameID = action.payload.id;
      },
      prepare: (id, history) => {
        return { payload: { id, history } };
      },
    },
    gameRetrieved(state, action) {
      state.loading = false;
      state.currentGame = action.payload.data;
      switch (state.currentGame.phase) {
        case `Deal`:
          // TODO leave numShuffles
          state.currentAction = {};
          break;
      }
    },
    exitGame: {
      reducer: state => {
        state.loading = false;
        state.currentGameID = '';
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
    refreshGame: {
      reducer: (state, action) => {
        if (state.currentGameID !== action.payload.id) {
          throw `bad game id: expected "${state.currentGameID}", got "${action.payload.id}"`;
        }
      },
      prepare: gameID => {
        return { payload: { id: gameID } };
      },
    },
    shuffleDeck(state) {
      isNaN(state.currentAction.numShuffles)
        ? (state.currentAction.numShuffles = 1)
        : (state.currentAction.numShuffles =
            state.currentAction.numShuffles + 1);
    },
    dealCards: {
      reducer: (state, action) => {
        // Nothing here?
      },
      prepare: history => {
        return { payload: { history } };
      },
    },
  },
});

export const { actions, reducer, name: sliceKey } = gameSlice;
