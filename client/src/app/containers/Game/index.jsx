import React from 'react';

import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import ActionBox from 'app/containers/Game/ActionBox';
import PlayerHand from 'app/containers/Game/PlayerHand';
import { gameSaga } from 'app/containers/Game/saga';
import RightSide from 'app/containers/Game/RightSide';
import { selectCurrentGame } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { selectCurrentUser } from 'auth/selectors';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const useStyles = makeStyles({
  gameArea: {
    width: '80%',
  },
});

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
  const classes = useStyles();
  const currentUser = useSelector(selectCurrentUser);
  const activeGame = useSelector(selectCurrentGame);

  return (
    <div>
      <Grid
        className={classes.gameArea}
        container
        xl
        spacing={1}
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
        >
          <Grid item sm container>
            <PlayerHand
              phase={activeGame.phase}
              hand={handForPlayer(activeGame, currentUser.id, 'across')}
            />
          </Grid>
          <Grid
            item
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
                isBlocking={Object.keys(activeGame.blocking_players).includes(
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
          <Grid item sm container>
            <PlayerHand
              mine
              phase={activeGame.phase}
              hand={activeGame.hands[currentUser.id]}
              pegged={activeGame.pegged_cards}
            />
          </Grid>
        </Grid>
      </Grid>
      <RightSide key='rightSidePanel' />
    </div>
  );
};

export default Game;
