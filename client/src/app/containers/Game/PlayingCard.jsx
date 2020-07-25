import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { makeStyles } from '@material-ui/core/styles';
import { red, grey } from '@material-ui/core/colors';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentAction } from './selectors';

const useStyles = makeStyles({
  root: {
    width: 96,
    height: 120,
  },
  unknown: {
    backgroundColor: grey[400],
  },
  value: {
    fontSize: 18,
  },
  redCard: {
    color: red[800],
  },
  blackCard: {
    color: grey[900],
  },
  suit: {
    fontSize: 36,
    justifyContent: 'center',
    alignItems: 'center',
    verticalAlign: 'center',
    textAlign: 'center',
  },
});

const PlayingCard = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const classes = useStyles();
  const dispatch = useDispatch();
  const currentAction = useSelector(selectCurrentAction);

  const useRed = !['Spades', 'Clubs'].includes(props.card.suit);

  if (!props.experimental) {
    if (props.card.name === 'unknown') {
      return <Card className={`${classes.root} ${classes.unknown}`}></Card>;
    }
    const suitEmojis = {
      Spades: '♠️',
      Clubs: '♣️',
      Diamonds: '♦️',
      Hearts: '♥️',
    };
    const valueStrings = {
      11: 'J',
      12: 'Q',
      13: 'K',
      1: 'A',
    };
    return (
      <Card className={classes.root}>
        <CardContent boxShadow={2}>
          <Typography
            variant='button'
            className={`${classes.value} ${
              useRed ? classes.redCard : classes.blackCard
            }`}
          >
            {valueStrings[props.card.value]
              ? valueStrings[props.card.value]
              : props.card.value}
          </Typography>
          <Typography
            className={`${classes.suit} ${
              useRed ? classes.redCard : classes.blackCard
            }`}
          >
            {suitEmojis[props.card.suit]}
          </Typography>
        </CardContent>
      </Card>
    );
  }

  if (!props.card) {
    return null;
  } else if (props.card.name === 'unknown') {
    // Currently, this returns a grayed out box, but it should show
    // a back of a card
    return (
      <div className='w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800' />
    );
  }

  let chosen = currentAction.selectedCards.indexOf(props.card) !== -1;
  let toggleChosen = () => {
    if (!props.disabled) {
      dispatch(actions.selectCard(props.card));
    }
  };

  return (
    <div
      onClick={props.mine ? toggleChosen : () => {}}
      className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
        props.disabled ? 'bg-gray-500' : 'bg-white'
      } ${useRed ? 'text-red-700' : 'text-black'}`}
      style={{
        position: 'relative',
        top: chosen ? '-10px' : '',
      }}
    >
      {props.card.name}
    </div>
  );
};

export default PlayingCard;
