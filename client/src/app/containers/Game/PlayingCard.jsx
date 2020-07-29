import React from 'react';

import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardMedia from '@material-ui/core/CardMedia';
import { makeStyles } from '@material-ui/core/styles';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const useStyles = makeStyles({
  root: {
    maxWidth: 96,
  },
  used: {
    opacity: 0.6,
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
  let value = valueStrings[card.value] ? valueStrings[card.value] : card.value;
  value += suitEmojis[card.suit];

  return (
    <Card
      variant={chosen ? 'outlined' : ''}
      onClick={mine ? toggleChosen : () => {}}
      className={`${classes.root} ${disabled ? classes.used : ''}`}
    >
      <CardActionArea disabled={!mine || disabled}>
        <CardMedia
          component='img'
          alt={value}
          image={`/cards/${
            card.name === 'unknown' ? 'background' : card.name
          }.svg`}
          title='Card'
        />
      </CardActionArea>
    </Card>
  );
};

PlayingCard.propTypes = {
  card: PropTypes.object.isRequired,
  disabled: PropTypes.bool.isRequired,
  experimental: PropTypes.bool.isRequired,
  mine: PropTypes.bool.isRequired,
};

export default PlayingCard;
