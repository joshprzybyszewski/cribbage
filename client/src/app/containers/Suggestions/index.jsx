import React from 'react';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import GridList from '@material-ui/core/GridList';
import SuggestionsTable from 'app/containers/Suggestions/SuggestionsTable';
import ChoosingCard from 'app/containers/Suggestions/ChoosingCard';

import { selectHandCards, getHandSuggestion } from 'app/containers/Suggestions/selectors';
import { suggestionsSaga } from 'app/containers/Suggestions/saga';
import { sliceKey, reducer, actions as sugActions, } from 'app/containers/Suggestions/slice';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const Suggestions = () => {
  useInjectSaga({ key: sliceKey, saga: suggestionsSaga });
  useInjectReducer({ key: sliceKey, reducer });
  const dispatch = useDispatch();

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
      <Button
        color='primary'
        variant='outlined'
        onClick={() => {
          dispatch(sugActions.getHandSuggestion());
        }}
      >
        Calculate
      </Button>
      <SuggestionsTable />
    </div>
  );
};

export default Suggestions;
