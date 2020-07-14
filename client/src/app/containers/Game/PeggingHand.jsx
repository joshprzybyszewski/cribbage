import React from 'react';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';

const PeggingHand = props => {
  if (!props.hand) {
    return null;
  }
  return (
    <Grid container justify='center' spacing={2} className='bg-green-800'>
      {props.hand.map(value => (
        <Grid key={value.name} item>
          <Paper
            className={{
              height: 140,
              width: 100,
            }}
          >
            {value.name}
          </Paper>
        </Grid>
      ))}
    </Grid>
  );
};

export default PeggingHand;
