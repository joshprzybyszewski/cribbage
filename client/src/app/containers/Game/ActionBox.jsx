import React from 'react';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';
import TextField from '@material-ui/core/TextField';

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
      ) : props.phase === 'Cut' ? (
        <Button
          disabled={!props.isBlocking}
          variant='contained'
          color='primary'
          endIcon={<CallSplitIcon />}
        >
          Cut
        </Button>
      ) : props.phase === 'Pegging' ? (
        <ButtonGroup
          orientation='vertical'
          color='primary'
          aria-label='vertical outlined primary button group'
        >
          <Button disabled={!props.isBlocking} color='secondary'>
            Say Go
          </Button>
          <Button
            disabled={!props.isBlocking}
            color='primary'
            endIcon={<SendIcon />}
          >
            Peg
          </Button>
        </ButtonGroup>
      ) : props.phase === 'Counting' ? (
        <form autoComplete='off'>
          <TextField
            id='handCountField'
            label='Hand Points'
            type='number'
            variant='outlined'
            InputLabelProps={{
              shrink: true,
            }}
          />
          <Button
            disabled={!props.isBlocking}
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
          >
            Count
          </Button>
        </form>
      ) : props.phase === 'CribCounting' ? (
        <form autoComplete='off'>
          <TextField
            id='cribCountField'
            label='Crib Points'
            type='number'
            variant='outlined'
            InputLabelProps={{
              shrink: true,
            }}
          />
          <Button
            disabled={!props.isBlocking}
            variant='contained'
            color='primary'
            endIcon={<SendIcon />}
          >
            Count
          </Button>
        </form>
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
