import React from 'react';

import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import SuggestionsTable from 'app/containers/Suggestions/SuggestionsTable';
import ChoosingCard from 'app/containers/Suggestions/ChoosingCard';

import { selectHandCards } from 'app/containers/Suggestions/selectors';
import { sliceKey, reducer } from 'app/containers/Suggestions/slice';
import { useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const Suggestions = () => {
  useInjectReducer({ key: sliceKey, reducer });

  const handCards = useSelector(selectHandCards);

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
      <SuggestionsTable/>
    </div>
  );
};

export default Suggestions;
