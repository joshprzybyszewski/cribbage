import React from 'react';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Grid from '@material-ui/core/Grid';

import { gameSaga } from './saga';
import { sliceKey, reducer } from './slice';
import DealAction from './DealAction';
import CribAction from './CribAction';
import CutAction from './CutAction';
import PegAction from './PegAction';
import CountHandAction from './CountHandAction';
import CountCribAction from './CountCribAction';

const ActionBox = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  return (
    <Grid item container justify='center' spacing={1}>
      {props.phase === 'Deal' ? (
        <DealAction isBlocking={props.isBlocking} />
      ) : props.phase === 'BuildCrib' ? (
        <CribAction isBlocking={props.isBlocking} />
      ) : props.phase === 'Cut' ? (
        <CutAction isBlocking={props.isBlocking} />
      ) : props.phase === 'Pegging' ? (
        <PegAction isBlocking={props.isBlocking} />
      ) : props.phase === 'Counting' ? (
        <CountHandAction isBlocking={props.isBlocking} />
      ) : props.phase === 'CribCounting' ? (
        <CountCribAction isBlocking={props.isBlocking} />
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
