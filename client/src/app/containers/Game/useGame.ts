import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

import { useAuth } from '../../../auth/useAuth';
import { RootState } from '../../../store/store';
import { useAlert } from '../Alert/useAlert';
import { Card, Game, Phase } from './models';
import {
    actions,
    CountCribAction,
    CountHandAction,
    CribAction,
    CutAction,
    DealAction,
    GameAction,
    PegAction,
} from './slice';
import { CreateGameResponse } from './types';

interface Result {
    game: Game;
    selectedCards: Card[];
    createGame: (playerIDs: string[]) => Promise<void>;
    loadActiveGame: (id: number) => Promise<void>;
    refreshGame: () => Promise<void>;
    toggleSelectedCard: (c: Card) => void;
    clearSelectedCards: () => void;
    submitDealAction: (a: DealAction) => Promise<void>;
    submitBuildCribAction: (a: CribAction) => Promise<void>;
    submitCutDeckAction: (a: CutAction) => Promise<void>;
    submitPegAction: (a: PegAction) => Promise<void>;
    submitCountHandAction: (a: CountHandAction) => Promise<void>;
    submitCountCribAction: (a: CountCribAction) => Promise<void>;
}

interface ServerCard {
    s: number;
    v: number;
}

interface ActionRequest {
    pID: string;
    gID: number;
    o: ReturnType<typeof mapPhaseToOverComes>;
    a: {
        ns?: number;
        cs?: ServerCard[];
        p?: number;
        sg?: boolean;
        c?: ServerCard;
        pts?: number;
    };
}

const mapPhaseToOverComes = (p: Phase) => {
    if (p === 'Deal') {
        return 0;
    }
    if (p === 'BuildCrib') {
        return 1;
    }
    if (p === 'Cut') {
        return 2;
    }
    if (p === 'Pegging') {
        return 3;
    }
    if (p === 'Counting') {
        return 4;
    }
    if (p === 'CribCounting') {
        return 5;
    }
    return -1;
};

const cardToServerCard = (c: Card): ServerCard => {
    const magicMap = {
        Spades: 0,
        Clubs: 1,
        Diamonds: 2,
        Hearts: 3,
    };
    return {
        s: magicMap[c.suit],
        v: c.value,
    };
};

function getPlayerAction<T extends GameAction>(
    myID: string,
    gID: number,
    phase: Phase,
    currentAction: T,
): ActionRequest {
    const request: ActionRequest = {
        pID: myID,
        gID,
        o: mapPhaseToOverComes(phase),
        a: {},
    };
    switch (phase) {
        case 'Deal':
            request.a = { ns: (currentAction as DealAction).numShuffles };
            break;
        case 'BuildCrib': {
            const cribAction = currentAction as CribAction;
            request.a = {
                cs: cribAction.selectedCards
                    ? cribAction.selectedCards.map(cardToServerCard)
                    : [],
            };
            break;
        }
        case 'Cut':
            request.a = { p: (currentAction as CutAction).percCut };
            break;
        case 'Pegging': {
            const pegAction = currentAction as PegAction;
            if (
                !pegAction.selectedCards ||
                pegAction.selectedCards.length !== 1
            ) {
                request.a = { sg: true };
            } else {
                request.a = { c: cardToServerCard(pegAction.selectedCards[0]) };
            }
            break;
        }
        case 'Counting':
            request.a = { pts: (currentAction as CountHandAction).points };
            break;
        case 'CribCounting':
            request.a = { pts: (currentAction as CountCribAction).points };
            break;
        default:
            request.a = {};
            break;
    }
    return request;
}

export function useGame(): Result {
    const gameState = useSelector((state: RootState) => state.game);
    const { currentUser } = useAuth();
    const { setAlert } = useAlert();
    const dispatch = useDispatch();

    const fetchGame = async (id: number) => {
        const response = await axios.get<Game>(
            `/game/${id}?player=${currentUser.id}`,
        );
        return response.data;
    };

    const loadActiveGame = async (id: number) => {
        try {
            const game = await fetchGame(id);
            dispatch(actions.setGame(game));
        } catch (err) {
            setAlert(err.response.data, 'error');
        }
    };

    const refreshGame = async () => {
        dispatch(actions.setLoading(true));
        try {
            const game = await fetchGame(gameState.currentGameID);
            dispatch(actions.setGame(game));
        } catch (err) {
            setAlert(err.response.data, 'error');
        }
        dispatch(actions.setLoading(false));
    };

    const createGame = async (playerIDs: string[]) => {
        dispatch(actions.setLoading(true));
        try {
            const createResult = await axios.post<CreateGameResponse>(
                `/create/game`,
                {
                    playerIDs,
                },
                // {
                //     baseURL: `lambda.hobbycribbage.com`,
                // },
            );
            const getResult = await axios.get<Game>(
                `/game/${createResult.data.id}`,
            );
            dispatch(actions.setGame(getResult.data));
        } catch (err) {
            setAlert(err.response.data, 'error');
        }
        dispatch(actions.setLoading(false));
    };

    const createActionHandler = (phase: Phase) => async (a: GameAction) => {
        dispatch(actions.setLoading(true));
        try {
            const request = getPlayerAction(
                currentUser.id,
                gameState.currentGameID,
                phase,
                a,
            );
            await axios.post('/action', request);
            await refreshGame();
        } catch (err) {
            setAlert(
                `handling action broke ${
                    err.response ? err.response.data : err
                }`,
                'error',
            );
        }
        dispatch(actions.setLoading(false));
    };

    return {
        game: gameState.currentGame,
        selectedCards: useSelector(
            (state: RootState) => state.game.selectedCards,
        ),
        createGame,
        loadActiveGame,
        refreshGame,
        toggleSelectedCard: (c: Card) =>
            dispatch(actions.toggleSelectedCard(c)),
        clearSelectedCards: () => dispatch(actions.clearSelectedCards()),
        submitDealAction: createActionHandler('Deal'),
        submitBuildCribAction: createActionHandler('BuildCrib'),
        submitCutDeckAction: createActionHandler('Cut'),
        submitPegAction: createActionHandler('Pegging'),
        submitCountHandAction: createActionHandler('Counting'),
        submitCountCribAction: createActionHandler('CribCounting'),
    };
}
