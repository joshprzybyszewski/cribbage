import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
  activeGamesPlayerID: '',
  activeGames: {},
};

const homeSlice = createSlice({
  name: 'home',
  initialState,
  reducers: {
    refreshActiveGames(state, action) {
      if (!action.payload.id) {
        // what should we do when refreshing with an ID we do not expect?
        throw `requires a playerID: got "${action.payload.id}"`;
      }
    },
    gotActiveGames(state, action) {
      state.activeGamesPlayerID = action.payload.player.id;
      state.activeGames = action.payload.activeGames;
    },
  },
});

export const { actions, reducer, name: sliceKey } = homeSlice;
