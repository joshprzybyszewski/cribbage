import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { AlertType } from './types';

export interface Alert {
    id: string;
    msg: string;
    type: AlertType;
}

export const initialState: Alert[] = [];

const alertSlice = createSlice({
    name: 'alert',
    initialState,
    reducers: {
        addAlert(state, action: PayloadAction<Alert>) {
            return [...state, action.payload];
        },
        removeAlert(state, action: PayloadAction<string>) {
            return state.filter(a => a.id !== action.payload);
        },
    },
});

export const { actions, reducer, name: sliceKey } = alertSlice;
