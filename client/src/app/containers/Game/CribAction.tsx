import React from 'react';

import Button from '@material-ui/core/Button';
import SendIcon from '@material-ui/icons/Send';

import { ActiveGame } from '../Home/slice';
import { ActionInputProps } from './types';
import { useGame } from './useGame';

const expNumCardsToCribForGame = (game: ActiveGame) => {
    if (game.players.length > 2) {
        return 1;
    }
    return 2;
};

const CribAction: React.FunctionComponent<ActionInputProps> = ({
    isBlocking,
}) => {
    const { game, selectedCards, submitBuildCribAction } = useGame();

    return (
        <Button
            disabled={
                !isBlocking ||
                selectedCards.length !==
                    expNumCardsToCribForGame(game.currentGame)
            }
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
            onClick={() => submitBuildCribAction({ selectedCards })}
        >
            Build Crib
        </Button>
    );
};

export default CribAction;
