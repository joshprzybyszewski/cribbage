import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';

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
    >
      {props.hand.map((card, index) => (
        <GridList key={card.name} item>
          <PlayingCard
            key={`handcard${index}`}
            card={card}
            mine={props.mine}
            disabled={
              props.phase === 'Pegging' &&
              props.pegged &&
              props.pegged.some(pc => pc.card.name === card.name)
            }
          />
        </GridList>
      ))}
    </Grid>
  );
};

export default PlayerHand;
