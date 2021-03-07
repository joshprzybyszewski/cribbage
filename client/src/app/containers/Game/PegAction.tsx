import React from 'react';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import SendIcon from '@material-ui/icons/Send';

import { ActionInputProps } from './types';
import { useGame } from './useGame';

const PegAction: React.FunctionComponent<ActionInputProps> = ({
    isBlocking,
}) => {
    const { selectedCards, submitPegAction } = useGame();

    return (
        <ButtonGroup
            orientation='vertical'
            color='primary'
            aria-label='vertical outlined primary button group'
        >
            <Button
                disabled={!isBlocking}
                color='secondary'
                // TODO it's probably semantically better if we also have a submitSayGoAction
                onClick={() => submitPegAction({ selectedCards: [] })}
            >
                Say Go
            </Button>
            <Button
                disabled={!isBlocking || selectedCards.length !== 1}
                color='primary'
                endIcon={<SendIcon />}
                onClick={() => submitPegAction({ selectedCards })}
            >
                Peg
            </Button>
        </ButtonGroup>
    );
};

export default PegAction;
