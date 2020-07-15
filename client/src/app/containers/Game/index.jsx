import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import Grid from '@material-ui/core/Grid';
import RefreshIcon from '@material-ui/icons/Refresh';
import IconButton from '@material-ui/core/IconButton';

import { selectCurrentUser } from '../../../auth/selectors';
import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentGame } from './selectors';
import ActionBox from './ActionBox';
import PlayingCard from './PlayingCard';
import PlayerHand from './PlayerHand';
import ScoreBoard from './ScoreBoard';

const showCutCard = phase => {
  return phase !== 'Deal' && phase !== 'BuildCrib' && phase !== 'Cut';
};

const handForPlayer = (game, myID, position) => {
  let isFourPlayer =
    game.teams.length === 2 && game.teams[0].players.length === 2;
  if (position === 'across') {
    if (game.teams.length === 3) {
      let secondPlayerID = game.teams.filter(
        t => !t.players.some(p => p.id === myID),
      )[1].players[0].id;
      return game.hands[secondPlayerID];
    } else if (isFourPlayer) {
      let partnerID = game.teams
        .filter(t => t.players.some(p => p.id === myID))[0]
        .players.filter(p => p.id !== myID)[0].id;
      return game.hands[partnerID];
    }
    let opponentID = game.teams.filter(
      t => !t.players.some(p => p.id === myID),
    )[0].players[0].id;
    return game.hands[opponentID];
  } else if (position === 'right') {
    if (isFourPlayer) {
      let rightID = game.teams
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
  let leftID = game.teams
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
    <Grid
      container
      xl
      spacing={1}
      className='bg-green-200'
      direction='row'
      justify='space-between'
    >
      <Grid
        item
        container
        md
        spacing={2}
        direction='column'
        align-content='space-between'
        className='bg-green-300'
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
            <div>{'TODO put an image of the deck here'}</div>,
            showCutCard(activeGame.phase)
              ? jsonCardToCard(activeGame.cut_card)
              : null,
          ]}
        </Grid>
      </Grid>
    </Grid>
  );
};

const jsonCardToCard = card => {
  return (
    <PlayingCard
      key={card.name}
      name={card.name}
      value={card.value}
      suit={card.suit}
    />
  );
};

export default Game;
