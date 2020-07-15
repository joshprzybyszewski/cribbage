import React from 'react';

import Grid from '@material-ui/core/Grid';

import PlayingCard from './PlayingCard';

const showOpponentsHand = phase => {
  return phase !== 'Deal';
};

const PlayerHand = props => {
  if (!props.hand || !showOpponentsHand(props.phase)) {
    return null;
  }
  return (
    <Grid
      item
      container
      direction={props.side ? 'column' : 'row'}
      justify='center'
      spacing={1}
      className='bg-green-800'
    >
      {props.hand.map((card, index) => (
        <Grid key={card.name} item>
          <PlayingCard
            key={`handcard${index}`}
            name={card.name}
            value={card.value}
            suit={card.suit}
            mine={props.mine}
            disabled={
              props.phase === 'Pegging' &&
              props.pegged &&
              props.pegged.some(pc => pc.card.name === card.name)
            }
          />
        </Grid>
      ))}
    </Grid>
  );
};

export default PlayerHand;
