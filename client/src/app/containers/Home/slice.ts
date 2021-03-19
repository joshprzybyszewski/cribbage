import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { User } from '../../../auth/slice';

interface Player extends User {
    color: string;
}

export interface ActiveGameResponse {
    player: User;
    activeGames: ActiveGame[];
}

export interface ActiveGame {
    gameID: number;
    players: Player[];
    created: Date;
    lastMove: Date;
}

export const initialState: ActiveGameResponse = {
    player: {
        id: '',
        name: '',
    },
    activeGames: [],
};

const homeSlice = createSlice({
    name: 'home',
    initialState,
    reducers: {
        setActiveGamesPlayerID(state, action: PayloadAction<string>) {
            return {
                ...state,
                player: { id: action.payload, name: '' },
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
