import { createSlice } from '@reduxjs/toolkit';

export const initialState = {};

const newGameSlice = createSlice({
    name: 'newGame',
    initialState,
    reducers: {
        createGame: {
            reducer: () => {
                // is there nothing to do here?
            },
            prepare: (opp1ID, opp2ID, teammateID, history) => {
                return { payload: { opp1ID, opp2ID, teammateID, history } };
            },
        },
    },
});

export const { actions, reducer, name: sliceKey } = newGameSlice;
