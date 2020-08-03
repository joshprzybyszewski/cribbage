import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';
import { useCurrentPlayerAndGame } from 'app/containers/Game/hooks';
import { actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';

const DealAction = ({ isBlocking }) => {
  const dispatch = useDispatch();
  const { currentUser, gameID } = useCurrentPlayerAndGame();
  const [numShuffles, setNumShuffles] = useState(0);

  return (
    <Grid item container spacing={2}>
      <Button
        disabled={!isBlocking}
        variant='contained'
        color='secondary'
        endIcon={<ShuffleIcon />}
        onClick={() => setNumShuffles(numShuffles + 1)}
      >
        Shuffle
      </Button>
      <Button
        disabled={!isBlocking || numShuffles <= 0}
        variant='contained'
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          setNumShuffles(0);
          dispatch(
            actions.dealCards({ userID: currentUser.id, gameID, numShuffles }),
          );
        }}
      >
        Deal
      </Button>
    </Grid>
  );
};

DealAction.propTypes = {
  isBlocking: PropTypes.bool.isRequired,
};

export default DealAction;
