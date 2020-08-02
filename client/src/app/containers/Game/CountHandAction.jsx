import React from 'react';

import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import SendIcon from '@material-ui/icons/Send';
import { useCurrentPlayerAndGame } from 'app/containers/Game/hooks';
import { actions } from 'app/containers/Game/slice';
import { useFormInput } from 'hooks/useFormInput';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';

const CountHandAction = ({ isBlocking, isCrib }) => {
  const [points, handlePointsChange] = useFormInput(0);
  const { currentUser, gameID } = useCurrentPlayerAndGame();
  const dispatch = useDispatch();
  return (
    <FormGroup row autoComplete='off'>
      <FormControl>
        <InputLabel htmlFor='component-simple'>
          {isCrib ? 'Crib' : 'Hand'} Points
        </InputLabel>
        <Input
          disabled={!isBlocking}
          id='component-simple'
          type='number'
          onChange={handlePointsChange}
        />
      </FormControl>
      <Button
        disabled={!isBlocking || points < 0}
        variant='contained'
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          dispatch(
            actions.countHand({
              userID: currentUser.id,
              gameID,
              points: Number(points),
              isCrib,
            }),
          );
        }}
      >
        Count
      </Button>
    </FormGroup>
  );
};

CountHandAction.propTypes = {
  isBlocking: PropTypes.bool.isRequired,
  isCrib: PropTypes.bool.isRequired,
};

export default CountHandAction;
