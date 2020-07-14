import React from 'react';

import Grid from '@material-ui/core/Grid';

import PlayingCard from './PlayingCard';

const PlayerHand = props => {
  if (!props.hand) {
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
      {props.hand.map(card => (
        <Grid key={card.name} item>
          <PlayingCard
            key={card.name}
            name={card.name}
            value={card.value}
            suit={card.suit}
          />
        </Grid>
      ))}
    </Grid>
  );
};

export default PlayerHand;
