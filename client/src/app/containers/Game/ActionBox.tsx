import React from 'react';

import Grid from '@material-ui/core/Grid';

import CountAction from './CountAction';
import CribAction from './CribAction';
import CutAction from './CutAction';
import DealAction from './DealAction';
import PegAction from './PegAction';
import { Phase } from './slice';

interface Props {
    phase: Phase;
    isBlocking: boolean;
}

const Action: React.FunctionComponent<Props> = ({ phase, isBlocking }) => {
    switch (phase) {
        case 'Deal':
            return <DealAction isBlocking={isBlocking} />;
        case 'BuildCrib':
            return <CribAction isBlocking={isBlocking} />;
        case 'Cut':
            return <CutAction isBlocking={isBlocking} />;
        case 'Pegging':
            return <PegAction isBlocking={isBlocking} />;
        case 'Counting':
            return <CountAction isBlocking={isBlocking} isCrib={false} />;
        case 'CribCounting':
            return <CountAction isBlocking={isBlocking} isCrib />;
        default:
            return null;
    }
};

const ActionBox: React.FunctionComponent<Props> = ({ phase, isBlocking }) => {
    return (
        <Grid item container justify='center' spacing={1}>
            <Action phase={phase} isBlocking={isBlocking} />
        </Grid>
    );
};

export default ActionBox;
