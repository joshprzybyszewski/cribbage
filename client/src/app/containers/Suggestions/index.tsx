import React from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import SuggestionsTable from './SuggestionsTable';
import ChoosingCard from './ChoosingCard';

import { useTossSuggestion } from './useTossSuggestion';

const Suggestions: React.FunctionComponent = () => {
  const { handCards, fetchSuggestions } = useTossSuggestion();

  return (
    <div className='fixed w-half-screen'>
      <Grid
        item
        container
        spacing={1}
      ><GridList>
          {handCards.map((card, index) => (
            <ChoosingCard
              key={`handcard${index}`}
              card={card}
            />
          ))}
        </GridList>
      </Grid>
      <Button
        color='primary'
        variant='outlined'
        onClick={fetchSuggestions}
      >
        Calculate
      </Button>
      <SuggestionsTable />
    </div>
  );
};

export default Suggestions;
