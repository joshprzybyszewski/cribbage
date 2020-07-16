import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';

import PlayingCard from './PlayingCard';

const CribHand = props => {
  if (!props.cards) {
    return null;
  }
  return (
    <Grid item container direction={'row'} justify='center' spacing={1}>
      <GridList item>
        {props.cards.map((card, index) => (
          <PlayingCard key={`cribcard${index}`} card={card} />
        ))}
      </GridList>
    </Grid>
  );
};

export default CribHand;
