import React, { useEffect, useState } from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import { selectSelectedCards } from 'app/containers/Game/selectors';
import { actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useDispatch, useSelector } from 'react-redux';

const useStyles = makeStyles({
  root: {
    width: 120,
    height: 160,
  },
  value: {
    fontSize: 14,
  },
  suit: {
    justifyContent: 'center',
    alignItems: 'center',
    verticalAlign: 'center',
    textAlign: 'center',
  },
});

const PlayingCard = ({ card, disabled, experimental, mine }) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [isChosen, setIsChosen] = useState(false);
  const selectedCards = useSelector(selectSelectedCards);

  useEffect(() => {
    if (selectedCards.map(c => c.name).includes(card.name)) {
      setIsChosen(true);
      return () => setIsChosen(false);
    }
  }, [selectedCards, card.name]);

  const isRed = !['Spades', 'Clubs'].includes(card.suit);

  if (experimental) {
    return (
      <Card className={classes.root}>
        <CardContent>
          <Typography
            className={classes.value}
            color={isRed ? 'red' : 'black'}
            gutterBottom
          >
            {card.value}
          </Typography>
          <Typography className={classes.suit}>
            {card.suit === 'Spades'
              ? '♠️'
              : card.suit === 'Clubs'
              ? '♣️'
              : card.suit === 'Diamonds'
              ? '♦️'
              : card.suit === 'Hearts'
              ? '♥️'
              : '?'}
          </Typography>
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

  const handleSelectCard = () => {
    if (!disabled) {
      if (isChosen) {
        dispatch(actions.unselectCard(card));
      } else {
        dispatch(actions.selectCard(card));
      }
    }
  };

  return (
    <div
      onClick={mine && handleSelectCard}
      className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
        disabled ? 'bg-gray-500' : 'bg-white'
      } ${isRed ? 'text-red-700' : 'text-black'}`}
      style={{
        position: 'relative',
        top: isChosen && '-10px',
      }}
    >
      {card.name}
    </div>
  );
};

PlayingCard.propTypes = {
  card: PropTypes.object.isRequired,
  disabled: PropTypes.bool,
  experimental: PropTypes.bool,
  mine: PropTypes.bool,
};

export default PlayingCard;
