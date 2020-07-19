import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import CallSplitIcon from '@material-ui/icons/CallSplit';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Slider from '@material-ui/core/Slider';
import SendIcon from '@material-ui/icons/Send';
import ShuffleIcon from '@material-ui/icons/Shuffle';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentAction, selectCurrentGame } from './selectors';

const expNumCardsToCribForGame = game => {
  if (game.teams.length === 3 || game.teams[0].players.length === 2) {
    return 1;
  }

  return 2;
};

const ActionBox = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);
  const activeGame = useSelector(selectCurrentGame);

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
      ) : props.phase === 'BuildCrib' ? (
        <Button
          disabled={
            !props.isBlocking ||
            currentAction.selectedCards.length !==
              expNumCardsToCribForGame(activeGame)
          }
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
            disabled={!props.isBlocking}
            orientation='vertical'
            getAriaValueText={perc}
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
      ) : props.phase === 'Pegging' ? (
        <ButtonGroup
          orientation='vertical'
          color='primary'
          aria-label='vertical outlined primary button group'
        >
          <Button
            disabled={!props.isBlocking}
            color='secondary'
            onClick={() => {
              dispatch(actions.pegCard());
            }}
          >
            Say Go
          </Button>
          <Button
            disabled={
              !props.isBlocking || currentAction.selectedCards.length !== 1
            }
            color='primary'
            endIcon={<SendIcon />}
            onClick={() => {
              dispatch(actions.pegCard());
            }}
          >
            Peg
          </Button>
        </ButtonGroup>
      ) : props.phase === 'Counting' ? (
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
      ) : props.phase === 'CribCounting' ? (
        <FormGroup row autoComplete='off'>
          <FormControl>
            <InputLabel htmlFor='component-simple'>Crib Points</InputLabel>
            <Input
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
              dispatch(actions.countCrib());
            }}
          >
            Count
          </Button>
        </FormGroup>
      ) : (
        'dev error!'
      )}
    </Grid>
  );
};

export default ActionBox;
