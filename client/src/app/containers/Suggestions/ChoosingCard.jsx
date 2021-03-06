import React from 'react';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import { sliceKey, reducer } from 'app/containers/Suggestions/slice';
import PropTypes from 'prop-types';
import { useInjectReducer } from 'redux-injectors';

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

function getSuit(card) {
  if (card.includes('C') || card.includes('c')) {
    return 'Clubs';
  } else if (card.includes('S') || card.includes('S')) {
    return 'Spades';
  } else if (card.includes('H') || card.includes('h')) {
    return 'Hearts';
  } else if (card.includes('D') || card.includes('d')) {
    return 'Diamonds';
  } 

  return 'Unknown'
}

function getValue(card) {
  if (card.startsWith('K') || card.startsWith('k') || card.startsWith('13')) {
    return 'K';
  } else if (card.startsWith('Q') || card.startsWith('q') || card.startsWith('12')) {
    return 'Q';
  } else if (card.startsWith('J') || card.startsWith('j') || card.startsWith('11')) {
    return 'J';
  } else if (card.startsWith('10')) {
    return '10';
  } else if (card.startsWith('A') || card.startsWith('a') || card.startsWith('1')) {
    return 'A';
  }

  return card.substr(0, card.length-1)
}

const ChoosingCard = ({ card, notEditable }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  const classes = useStyles();

  const suit = getSuit(card);
  const value = getValue(card);
  const useRed = !['Spades', 'Clubs'].includes(suit);

  const updateValue=(_) => {
    // TODO figure out how to get scroll capturing to work
    // TODO update the value of this card in the state
    console.log(`scrolled: ${card}`);
  };
  const updateSuit=(_) => {
    // TODO update the state so that this card increments suits
    console.log(`clicked: ${card}`);
  };

    return (
      <Card 
      onScroll={updateValue}
      onClick={updateSuit}
      className={classes.root}>
        <CardContent>
          <Typography
            className={classes.value}
            color={useRed ? 'red' : 'black'}
            gutterBottom
          >
            {value}
          </Typography>
          <Typography className={classes.suit}>
            {suit === 'Spades'
              ? '♠️'
              : suit === 'Clubs'
              ? '♣️'
              : suit === 'Diamonds'
              ? '♦️'
              : suit === 'Hearts'
              ? '♥️'
              : '?'}
          </Typography>
        </CardContent>
      </Card>
    );
};

ChoosingCard.propTypes = {
  card: PropTypes.string.isRequired,
  notEditable: PropTypes.bool,
};

export default ChoosingCard;
