import React from 'react';

import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import RefreshIcon from '@material-ui/icons/Refresh';
import ActionBox from 'app/containers/Game/ActionBox';
import CribHand from 'app/containers/Game/CribHand';
import PlayerHand from 'app/containers/Game/PlayerHand';
import PlayingCard from 'app/containers/Game/PlayingCard';
import { gameSaga } from 'app/containers/Game/saga';
import ScoreBoard from 'app/containers/Game/ScoreBoard';
import { selectCurrentGame } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { selectCurrentUser } from 'auth/selectors';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const showCutCard = phase => {
  return !['Deal', 'BuildCrib', 'Cut'].includes(phase);
};

const handForPlayer = (game, myID, position) => {
  const isFourPlayer =
    game.teams.length === 2 && game.teams[0].players.length === 2;
  if (position === 'across') {
    if (game.teams.length === 3) {
      const secondPlayerID = game.teams.filter(
        t => !t.players.some(p => p.id === myID),
      )[1].players[0].id;
      return game.hands[secondPlayerID];
    } else if (isFourPlayer) {
      const partnerID = game.teams
        .filter(t => t.players.some(p => p.id === myID))[0]
        .players.filter(p => p.id !== myID)[0].id;
      return game.hands[partnerID];
    }
    const opponentID = game.teams.filter(
      t => !t.players.some(p => p.id === myID),
    )[0].players[0].id;
    return game.hands[opponentID];
  } else if (position === 'right') {
    if (isFourPlayer) {
      const rightID = game.teams
        .filter(t => t.players.some(p => p.id !== myID))[0]
        .players.filter(p => p.id !== myID)[1].id;
      return game.hands[rightID];
    }
    // nothing!
    return null;
  } else if (position !== 'left' || !isFourPlayer) {
    return null;
  }
  // position is left
  const leftID = game.teams
    .filter(t => t.players.some(p => p.id !== myID))[0]
    .players.filter(p => p.id !== myID)[0].id;
  return game.hands[leftID];
};

const Game = () => {
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
    <Grid container xl spacing={1} direction='row' justify='space-between'>
      <Grid
        item
        container
        md
        spacing={2}
        direction='column'
        align-content='space-between'
      >
        <Grid item xs sm container>
          <PlayerHand
            phase={activeGame.phase}
            hand={handForPlayer(activeGame, currentUser.id, 'across')}
          />
        </Grid>
        <Grid
          item
          xs
          md
          container
          justify='space-between'
          align-content='center'
        >
          <Grid item>
            <PlayerHand
              side
              phase={activeGame.phase}
              hand={handForPlayer(activeGame, currentUser.id, 'left')}
            />
          </Grid>
          <Grid item>
            <ActionBox
              phase={activeGame.phase}
              isBlocking={activeGame.blocking_players.hasOwnProperty(
                currentUser.id,
              )}
            />
          </Grid>
          <Grid item>
            <PlayerHand
              side
              phase={activeGame.phase}
              hand={handForPlayer(activeGame, currentUser.id, 'right')}
            />
          </Grid>
        </Grid>
        <Grid item xs sm container>
          <PlayerHand
            mine
            phase={activeGame.phase}
            hand={activeGame.hands[currentUser.id]}
            pegged={activeGame.pegged_cards}
          />
        </Grid>
      </Grid>
      <Grid item container xs direction='column' spacing={1}>
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
            showCutCard(activeGame.phase) ? (
              <PlayingCard key='cutCard' card={activeGame.cut_card} />
            ) : (
              <div key='deckTODOdiv'>
                {'TODO put an image of the deck here'}
              </div>
            ),
            <CribHand cards={activeGame.crib} />,
            <div key='currentPeg'>
              {activeGame.phase === 'Pegging'
                ? `Current Peg: ${
                    activeGame.current_peg ? activeGame.current_peg : 0
                  }`
                : ''}
            </div>,
          ]}
        </Grid>
      </Grid>
    </Grid>
  );
};

export default Game;
