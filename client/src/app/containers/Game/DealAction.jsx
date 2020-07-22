import React from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const DealAction = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);

  return (
    <Grid item container spacing={2}>
      <Button
        disabled={!props.isBlocking}
        variant='contained'
        color='secondary'
        endIcon={<ShuffleIcon />}
        onClick={() => {
          dispatch(actions.shuffleDeck());
        }}
      >
        Shuffle
      </Button>
      <Button
        disabled={!props.isBlocking || currentAction.numShuffles <= 0}
        variant='contained'
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          dispatch(actions.dealCards());
        }}
      >
        Deal
      </Button>
    </Grid>
  );
};

export default DealAction;
