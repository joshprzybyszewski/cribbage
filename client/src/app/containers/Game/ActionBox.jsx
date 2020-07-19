import React from 'react';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Grid from '@material-ui/core/Grid';

import { gameSaga } from './saga';
import { sliceKey, reducer } from './slice';
import { DealAction } from './DealAction';
import { CribAction } from './CribAction';
import { CutAction } from './CutAction';
import { PegAction } from './PegAction';
import { CountHandAction } from './CountHandAction';
import { CountCribAction } from './CountCribAction';

const ActionBox = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  return (
    <Grid item container justify='center' spacing={1}>
      {props.phase === 'Deal' ? (
        <DealAction />
      ) : props.phase === 'BuildCrib' ? (
        <CribAction />
      ) : props.phase === 'Cut' ? (
        <CutAction />
      ) : props.phase === 'Pegging' ? (
        <PegAction />
      ) : props.phase === 'Counting' ? (
        <CountHandAction />
      ) : props.phase === 'CribCounting' ? (
        <CountCribAction />
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
