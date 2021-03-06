import { nanoid } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';

export const initialState = {
    alerts: [],
};

const alertSlice = createSlice({
    name: 'alert',
    initialState,
    reducers: {
        addAlert: {
            reducer: (state, action) => {
                state.alerts = [...state.alerts, action.payload];
            },
            prepare: (msg, type) => {
                const id = nanoid();
                return { payload: { id, msg, type } };
            },
        },
        removeAlert(state, action) {
            state.alerts = state.alerts.filter(a => a.id !== action.payload);
        },
    },
});

export const { actions, reducer, name: sliceKey } = alertSlice;
