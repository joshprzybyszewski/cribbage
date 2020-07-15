import React from 'react';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';

const ActionBox = props => {
  return (
    <Grid item container justify='center' spacing={1}>
      {props.phase === 'Deal' ? (
        <Grid item container spacing={2}>
          <Button
            disabled={!props.isBlocking}
            variant='contained'
            color='secondary'
            endIcon={<ShuffleIcon />}
          >
            Shuffle
          </Button>
          <Button
            disabled={!props.isBlocking}
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
          >
            Deal
          </Button>
        </Grid>
      ) : props.phase === 'BuildCrib' ? (
        <Button
          disabled={!props.isBlocking}
          variant='contained'
          color='primary'
          endIcon={<SendIcon />}
        >
          Build Crib
        </Button>
      ) : props.phase === 'Pegging' ? (
        <ButtonGroup
          disabled={!props.isBlocking}
          orientation='vertical'
          color='primary'
          aria-label='vertical contained primary button group'
          variant='text'
        >
          <Button
            disabled={!props.isBlocking}
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
          >
            Peg
          </Button>
          {/* <Button>Peg</Button> */}
          <Button>Say Go</Button>
        </ButtonGroup>
      ) : props.phase === 'Counting' ? (
        // TODO add an input, disable if I'm not blocking
        <Button disabled={!props.isBlocking}>Count Hand</Button>
      ) : props.phase === 'CribCounting' ? (
        // TODO add an input, disable if I'm not blocking
        <Button disabled={!props.isBlocking}>Count Crib</Button>
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
