import React from 'react';

import { red, grey } from '@material-ui/core/colors';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const useStyles = makeStyles({
  root: {
    width: 96,
    height: 120,
    display: 'flex',
    flex: '1 0 auto',
  },
  content: {
    flex: '1 0 auto',
  },
  unknown: {
    backgroundColor: grey[600],
  },
  used: {
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

const PlayingCard = ({ card, disabled, experimental, mine }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const classes = useStyles();
  const dispatch = useDispatch();
  const currentAction = useSelector(selectCurrentAction);

  const useRed = !['Spades', 'Clubs'].includes(card.suit);

  const chosen = currentAction.selectedCards.indexOf(card) !== -1;
  const toggleChosen = () => {
    if (!disabled) {
      dispatch(actions.selectCard(card));
    }
  };

  if (!experimental) {
    if (card.name === 'unknown') {
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
    let value = valueStrings[card.value]
      ? valueStrings[card.value]
      : card.value;
    value += suitEmojis[card.suit];
    return (
      <Card
        variant={chosen ? 'outlined' : ''}
        onClick={mine ? toggleChosen : () => {}}
        className={`${classes.root} ${disabled ? classes.used : ''}`}
      >
        <CardContent boxShadow={2} className={classes.content}>
          <CardActionArea disabled={disabled || !mine}>
            <Typography
              variant='button'
              className={`${classes.value} ${
                useRed ? classes.redCard : classes.blackCard
              }`}
            >
              {value}
            </Typography>
            <Typography
              className={`${classes.suit} ${
                useRed ? classes.redCard : classes.blackCard
              }`}
            >
              {suitEmojis[card.suit]}
            </Typography>
            <Typography
              variant='button'
              className={`${classes.value} ${
                useRed ? classes.redCard : classes.blackCard
              }`}
            >
              {value}
            </Typography>
          </CardActionArea>
        </CardContent>
      </Card>
    );
  }

  if (!card) {
    return null;
  } else if (card.name === 'unknown') {
    // Currently, this returns a grayed out box, but it should show
    // a back of a card
    return (
      <div className='w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800' />
    );
  }

  return (
    <div
      onClick={mine ? toggleChosen : () => {}}
      className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
        disabled ? 'bg-gray-500' : 'bg-white'
      } ${useRed ? 'text-red-700' : 'text-black'}`}
      style={{
        position: 'relative',
        top: chosen ? '-10px' : '',
      }}
    >
      {card.name}
    </div>
  );
};

PlayingCard.propTypes = {
  card: PropTypes.object.isRequired,
  disabled: PropTypes.bool.isRequired,
  experimental: PropTypes.bool.isRequired,
  mine: PropTypes.bool.isRequired,
};

export default PlayingCard;
