import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { User } from '../../../auth/slice';

interface Player extends User {
    color: string;
}

export interface ActiveGame {
    gameID: string;
    players: Player[];
    created: Date;
    lastMove: Date;
}

interface HomeState {
    activeGamesPlayerID: string;
    activeGames: ActiveGame[];
}

export const initialState: HomeState = {
    activeGamesPlayerID: '',
    activeGames: [],
};

const homeSlice = createSlice({
    name: 'home',
    initialState,
    reducers: {
        setActiveGamesPlayerID(state, action: PayloadAction<string>) {
            return {
                ...state,
                activeGamesPlayerID: action.payload,
            };
        },
        setActiveGames(state, action: PayloadAction<ActiveGame[]>) {
            return {
                ...state,
                activeGames: action.payload,
            };
        },
    },
});

export const { actions, reducer, name: sliceKey } = homeSlice;
