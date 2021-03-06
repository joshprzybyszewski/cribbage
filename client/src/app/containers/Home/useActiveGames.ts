import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

import { User } from '../../../auth/slice';
import { useAuth } from '../../../auth/useAuth';
import { RootState } from '../../../store/store';
import { useAlert } from '../Alert/useAlert';
import { actions, ActiveGame } from './slice';

interface ReturnType {
    games: ActiveGame[];
    refreshGames: () => Promise<void>;
}

interface ActiveGamesResponse {
    activeGames: ActiveGame[];
    player: User;
}

export function useActiveGames(): ReturnType {
    const { currentUser } = useAuth();
    const { setAlert } = useAlert();
    const dispatch = useDispatch();
    return {
        games: useSelector((state: RootState) => state.home.activeGames),
        refreshGames: async () => {
            if (currentUser.id.length > 0) {
                try {
                    const res = await axios.get<ActiveGamesResponse>(
                        `/games/active?playerID=${currentUser.id}`,
                    );
                    dispatch(
                        actions.setActiveGamesPlayerID(res.data.player.id),
                    );
                    dispatch(actions.setActiveGames(res.data.activeGames));
                } catch (err) {
                    dispatch(setAlert(err.response.data, 'error'));
                }
            }
        },
    };
}
