import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  currentGameID: '',
  currentGame: {},
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
  },
});

export const { actions, reducer, name: sliceKey } = gameSlice;
