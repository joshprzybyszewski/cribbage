import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import PlayingCard from 'app/containers/Game/PlayingCard';
import PropTypes from 'prop-types';

const CribHand = ({ cards }) => {
  return !cards ? null : (
    <Grid item container direction={'row'} justify='center' spacing={1}>
      <GridList>
        {cards.map((card, index) => (
          <PlayingCard key={`cribcard${index}`} card={card} />
        ))}
      </GridList>
    </Grid>
  );
};

CribHand.propTypes = {
  cards: PropTypes.array,
};

export default CribHand;
