import React from 'react';

import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import SendIcon from '@material-ui/icons/Send';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const CountHandAction = ({ isBlocking }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);

  return (
    <FormGroup row autoComplete='off'>
      <FormControl>
        <InputLabel htmlFor='component-simple'>Hand Points</InputLabel>
        <Input
          disabled={!isBlocking}
          id='component-simple'
          type='number'
          onChange={event => {
            dispatch(actions.claimPoints(Number(event.target.value)));
          }}
        />
      </FormControl>
      <Button
        disabled={!isBlocking || currentAction.points < 0}
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

CountHandAction.propTypes = {
  isBlocking: PropTypes.bool.isRequired,
};

export default CountHandAction;
