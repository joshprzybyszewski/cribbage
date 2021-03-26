import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';

import { ActionInputProps } from './types';
import { useGame } from './useGame';

const DealAction: React.FunctionComponent<ActionInputProps> = ({
    isBlocking,
}) => {
    const { submitDealAction } = useGame();
    const [numShuffles, setNumShuffles] = useState(0);
    return (
        <Grid item container spacing={2}>
            <Button
                disabled={!isBlocking}
                variant='contained'
                color='secondary'
                endIcon={<ShuffleIcon />}
                onClick={() => setNumShuffles(prev => prev + 1)}
            >
                Shuffle
            </Button>
            <Button
                disabled={!isBlocking || numShuffles <= 0}
                variant='contained'
                color='primary'
                endIcon={<SendIcon />}
                onClick={() => submitDealAction({ numShuffles })}
            >
                Deal
            </Button>
        </Grid>
    );
};

export default DealAction;
