import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '../../../store/store';
import { useAlert } from '../Alert/useAlert';
import { Card } from '../Game/models';
import {
    actions,
    TossSuggestion,
} from './slice';


interface Result {
    handCards: Card[];
    suggestedHands: TossSuggestion[];
    isLoading: boolean;
    fetchSuggestions: () => Promise<void>;
    updateCard: (prev: Card, cur: Card) => void;
}

export function useTossSuggestion(): Result {
    const suggestionsState = useSelector((state: RootState) => state.suggestions);
    const { setAlert } = useAlert();
    const dispatch = useDispatch();

    const fetchSuggestions = async () => {
        dispatch(actions.setLoading(true));
        const currentHand = '1c,2c,3c,4c,5c,6c';
        try {
            const getResult = await axios.get<TossSuggestion[]>(
                `/suggest/hand?dealt=${currentHand}`
            );
            dispatch(actions.setSuggestionResult(getResult.data));
        } catch (err) {
            setAlert(err.response.data, 'error');
        }
        dispatch(actions.setLoading(false));
    };

    return {
        handCards: suggestionsState.handCards,
        isLoading: suggestionsState.loading,
        suggestedHands: suggestionsState.suggestedHands,
        fetchSuggestions,
        updateCard: (p: Card, c: Card) =>
            dispatch(actions.updateCard({prev: p, cur: c})),
    };
}