import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import SendIcon from '@material-ui/icons/Send';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentAction } from './selectors';

const CountHandAction = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);

  return (
    <FormGroup row autoComplete='off'>
      <FormControl>
        <InputLabel htmlFor='component-simple'>Hand Points</InputLabel>
        <Input
          disabled={!props.isBlocking}
          id='component-simple'
          type='number'
          onChange={event => {
            dispatch(actions.claimPoints(Number(event.target.value)));
          }}
        />
      </FormControl>
      <Button
        disabled={!props.isBlocking || currentAction.points < 0}
        variant='contained'
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          dispatch(actions.countHand());
        }}
      >
        Count
      </Button>
    </FormGroup>
  );
};

export default CountHandAction;
