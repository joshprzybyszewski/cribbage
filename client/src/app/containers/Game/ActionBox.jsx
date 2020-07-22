import React from 'react';

import Grid from '@material-ui/core/Grid';
import CountCribAction from 'app/containers/Game/CountCribAction';
import CountHandAction from 'app/containers/Game/CountHandAction';
import CribAction from 'app/containers/Game/CribAction';
import CutAction from 'app/containers/Game/CutAction';
import DealAction from 'app/containers/Game/DealAction';
import PegAction from 'app/containers/Game/PegAction';
import { gameSaga } from 'app/containers/Game/saga';
import { sliceKey, reducer } from 'app/containers/Game/slice';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

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
