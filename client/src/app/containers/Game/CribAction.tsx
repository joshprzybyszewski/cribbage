import React from 'react';

import Button from '@material-ui/core/Button';
import SendIcon from '@material-ui/icons/Send';

import { Game } from './models';
import { ActionInputProps } from './types';
import { useGame } from './useGame';

const expNumCardsToCribForGame = (game: Game) => {
    const numPlayers = game.teams.reduce(
        (num, team) => team.players.length + num,
        0,
    );
    if (numPlayers > 2) {
        return 1;
    }
    return 2;
};

const CribAction: React.FunctionComponent<ActionInputProps> = ({
    isBlocking,
}) => {
    const {
        clearSelectedCards,
        game,
        selectedCards,
        submitBuildCribAction,
    } = useGame();

    const handleClick = async () => {
        // TODO the submit build crib action should just take care of this probably
        await submitBuildCribAction({ selectedCards });
        clearSelectedCards();
    };

    return (
        <Button
            disabled={
                !isBlocking ||
                selectedCards.length !== expNumCardsToCribForGame(game)
            }
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
            onClick={handleClick}
        >
            Build Crib
        </Button>
    );
};

export default CribAction;
