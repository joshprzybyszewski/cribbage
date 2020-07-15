import React from 'react';

import Grid from '@material-ui/core/Grid';

import PlayingCard from './PlayingCard';

const CribHand = props => {
  if (!props.cards) {
    return null;
  }
  return (
    <Grid item container direction={'row'} justify='center' spacing={1}>
      {props.cards.map((card, index) => (
        <Grid key={card.name} item>
          <PlayingCard
            key={`cribcard${index}`}
            name={card.name}
            value={card.value}
            suit={card.suit}
            mine={props.mine}
          />
        </Grid>
      ))}
    </Grid>
  );
};

export default CribHand;
