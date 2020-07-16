import React from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import Grid from '@material-ui/core/Grid';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import Slider from '@material-ui/core/Slider';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';

const ActionBox = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  function perc(value) {
    return `${value}%`;
  }

  return (
    <Grid item container justify='center' spacing={1}>
      {props.phase === 'Deal' ? (
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
            disabled={!props.isBlocking}
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
      ) : props.phase === 'BuildCrib' ? (
        <Button
          disabled={!props.isBlocking}
          variant='contained'
          color='primary'
          endIcon={<SendIcon />}
          onClick={() => {
            dispatch(actions.buildCrib());
          }}
        >
          Build Crib
        </Button>
      ) : props.phase === 'Cut' ? (
        <div>
          <Slider
            orientation='vertical'
            getAriaValueText={perc}
            defaultValue={50}
            aria-labelledby='vertical-slider'
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
