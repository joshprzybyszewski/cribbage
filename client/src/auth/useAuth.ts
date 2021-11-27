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
    const baseURL = `lambda.hobbycribbage.com`;

    return {
        currentUser,
        isLoggedIn,
        login: async (id: string) => {
            dispatch(actions.setLoading(true));
            try {
                const res = await axios.get<UserResponse>(
                    `/player/${id}`,
                    {
                        baseURL,
                    },
                    );
                dispatch(actions.setUser(res.data.player));
            } catch (err) {
                dispatch(actions.clearUser());
                setAlert(err.response.data, 'error');
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
                    '/create/player',
                    request,
                    {
                        baseURL,
                    },
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
