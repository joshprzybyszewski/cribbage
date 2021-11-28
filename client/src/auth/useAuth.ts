import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

import { useAlert } from '../app/containers/Alert/useAlert';
import { RootState } from '../store/store';
import { actions, User } from './slice';

interface ReturnType {
    currentUser: User;
    isLoggedIn: boolean;
    login: (id: string) => Promise<void>;
    logout: () => void;
    register: (name: string, id: string) => Promise<void>;
}

interface UserResponse {
    player: User;
}

interface RegisterRequest {
    player: User;
}

export function useAuth(): ReturnType {
    const { setAlert } = useAlert();
    const dispatch = useDispatch();
    const { currentUser, isLoggedIn } = useSelector(
        (state: RootState) => state.auth,
    );
    const base = `https://lambda.hobbycribbage.com`;

    return {
        currentUser,
        isLoggedIn,
        login: async (id: string) => {
            dispatch(actions.setLoading(true));
            try {
                const res = await axios.get<UserResponse>(
                    `${base}/player/${id}`,
                );
                dispatch(actions.setUser(res.data.player));
            } catch (err) {
                dispatch(actions.clearUser());
                if (err) {
                    if (err.response) {
                        setAlert(err.response.data, 'error');
                    } else {
                        setAlert(`no err.response ${err}`, 'error');
                    }
                } else {
                    setAlert('no err', 'error');
                }
            }
            dispatch(actions.setLoading(false));
        },
        logout: () => dispatch(actions.clearUser()),
        register: async (name: string, id: string) => {
            const request: RegisterRequest = {
                player: {
                    id,
                    name,
                },
            };
            dispatch(actions.setLoading(true));
            try {
                const res = await axios.post<UserResponse>(
                    `${base}/create/player`,
                    request,
                );
                dispatch(actions.setUser(res.data.player));
                setAlert('Registration successful!', 'success');
            } catch (err) {
                dispatch(actions.clearUser());
                setAlert(err.response.data, 'error');
            }
            dispatch(actions.setLoading(false));
        },
    };
}
