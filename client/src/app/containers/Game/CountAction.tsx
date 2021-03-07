import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import SendIcon from '@material-ui/icons/Send';

import { ActionInputProps } from './types';
import { useGame } from './useGame';

interface Props extends ActionInputProps {
    isCrib: boolean;
}

const CountAction: React.FunctionComponent<Props> = ({
    isBlocking,
    isCrib,
}) => {
    const [points, setPoints] = useState(0);
    const { submitCountHandAction, submitCountCribAction } = useGame();

    const submitAction = isCrib ? submitCountCribAction : submitCountHandAction;
    const componentId = isCrib ? 'count-crib-input' : 'count-hand-input';

    return (
        <FormGroup row>
            <FormControl>
                <InputLabel htmlFor={componentId}>
                    {isCrib ? 'Crib' : 'Hand'} Points
                </InputLabel>
                <Input
                    disabled={!isCrib && !isBlocking}
                    id={componentId}
                    type='number'
                    onChange={event => {
                        setPoints(Number(event.target.value));
                    }}
                />
            </FormControl>
            <Button
                disabled={!isBlocking || points < 0}
                variant='contained'
                color='primary'
                endIcon={<SendIcon />}
                onClick={() => submitAction({ points })}
            >
                Count
            </Button>
        </FormGroup>
    );
};

export default CountAction;
