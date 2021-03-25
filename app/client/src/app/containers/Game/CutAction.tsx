import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Slider from '@material-ui/core/Slider';
import CallSplitIcon from '@material-ui/icons/CallSplit';

import { ActionInputProps } from './types';
import { useGame } from './useGame';

const useStyles = makeStyles({
    container: {
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
    },
});

const CutAction: React.FunctionComponent<ActionInputProps> = ({
    isBlocking,
}) => {
    const [percentage, setPercentage] = useState(50);
    const { submitCutDeckAction } = useGame();
    const classes = useStyles();

    return (
        <div className={classes.container}>
            <Slider
                disabled={!isBlocking}
                orientation='vertical'
                getAriaValueText={value => `${value}%`}
                value={percentage}
                min={0}
                max={100}
                aria-labelledby='vertical-slider'
                onChange={(_, newValue) => setPercentage(newValue as number)}
            />
            <Button
                disabled={!isBlocking}
                variant='contained'
                color='primary'
                endIcon={<CallSplitIcon />}
                onClick={() =>
                    submitCutDeckAction({ percCut: percentage / 100 })
                }
            >
                Cut
            </Button>
        </div>
    );
};

export default CutAction;
