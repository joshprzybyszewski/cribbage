import React from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import { nanoid } from '@reduxjs/toolkit';

import ChoosingCard from './ChoosingCard';
import SuggestionsTable from './SuggestionsTable';
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
          {handCards.map((card) => (
            <ChoosingCard
              key={`handcard${card.name}${nanoid()}`}
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
