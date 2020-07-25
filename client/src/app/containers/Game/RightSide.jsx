import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Grid from '@material-ui/core/Grid';
import RefreshIcon from '@material-ui/icons/Refresh';
import IconButton from '@material-ui/core/IconButton';
import Container from '@material-ui/core/Container';

import { selectCurrentUser } from '../../../auth/selectors';
import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentGame } from './selectors';

import ScoreBoard from './ScoreBoard';
import PlayingCard from './PlayingCard';
import CribHand from './CribHand';

const showCutCard = phase => {
  return !['Deal', 'BuildCrib', 'Cut'].includes(phase);
};

const RightSide = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const currentUser = useSelector(selectCurrentUser);
  const activeGame = useSelector(selectCurrentGame);

  // event handlers
  const onRefreshCurrentGame = id => {
    dispatch(actions.refreshGame(id, history));
  };

  return (
    <Container container xs spacing={1}>
      <Grid item>
        <IconButton
          aria-label='refresh'
          onClick={() => onRefreshCurrentGame(activeGame.id)}
        >
          <RefreshIcon />
        </IconButton>
        <ScoreBoard
          teams={activeGame.teams}
          current_dealer={activeGame.current_dealer}
        />
      </Grid>
      <Grid item>
        {[
          <div key='currentPeg'>
            {activeGame.phase === 'Pegging'
              ? `Current Peg: ${
                  activeGame.current_peg ? activeGame.current_peg : 0
                }`
              : ''}
          </div>,
          showCutCard(activeGame.phase) ? (
            <PlayingCard key='cutCard' card={activeGame.cut_card} />
          ) : (
            <div key='deckTODOdiv'>{'TODO put an image of the deck here'}</div>
          ),
          <CribHand cards={activeGame.crib} />,
        ]}
      </Grid>
    </Container>
  );
};

export default RightSide;
