import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import PlayingCard from 'app/containers/Game/PlayingCard';

const CribHand = props => {
  return (
    !props.cards || (
      <Grid item container direction={'row'} justify='center' spacing={1}>
        <GridList>
          {props.cards.map((card, index) => (
            <PlayingCard key={`cribcard${index}`} card={card} />
          ))}
        </GridList>
      </Grid>
    )
  );
};

export default CribHand;
