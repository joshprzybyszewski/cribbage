import React from 'react';

import Button from '@material-ui/core/Button';
import Slider from '@material-ui/core/Slider';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import { gameSaga } from 'app/containers/Game/saga';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const CutAction = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  return (
    <div>
      <Slider
        disabled={!props.isBlocking}
        orientation='vertical'
        getAriaValueText={value => {
          return `${value}%`;
        }}
        defaultValue={50}
        aria-labelledby='vertical-slider'
        onChange={event => {
          dispatch(actions.claimPoints(Number(event.target.value) / 100));
        }}
      />
      <Button
        disabled={!props.isBlocking}
        variant='contained'
        color='primary'
        endIcon={<CallSplitIcon />}
        onClick={() => {
          // TODO get the value of the Slider and use that to cut
          dispatch(actions.cutDeck());
        }}
      >
        Cut
      </Button>
    </div>
  );
};

export default CutAction;
