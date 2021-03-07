import { nanoid } from '@reduxjs/toolkit';
import { useDispatch } from 'react-redux';

import { actions } from './slice';
import { AlertType } from './types';

export interface AlertSetter {
    (message: string, type: AlertType): string;
}

interface ReturnType {
    setAlert: AlertSetter;
    removeAlert: (id: string) => void;
}

export function useAlert(timeout = 5000): ReturnType {
    const dispatch = useDispatch();

    return {
        setAlert: (message: string, type: AlertType) => {
            const id = nanoid();
            dispatch(actions.addAlert({ id, msg: message, type }));
            setTimeout(() => {
                dispatch(actions.removeAlert(id));
            }, timeout);
            return id;
        },
        removeAlert: (id: string) => {
            dispatch(actions.removeAlert(id));
        },
    };
}
