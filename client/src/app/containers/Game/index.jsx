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
import PlayingCard from './PlayingCard';
import PlayerHand from './PlayerHand';
import ScoreBoard from './ScoreBoard';

const showCutCard = phase => {
  return phase !== 'Deal' && phase !== 'BuildCrib' && phase !== 'Cut';
};

const handForPlayer = (game, myID, position) => {
  if (position === 'across') {
    if (game.teams.length === 3) {
      let secondPlayerID = game.teams.filter(
        t => !t.players.some(p => p.id === myID),
      )[1].players[0].id;
      return game.hands[secondPlayerID];
    } else if (game.teams.length === 4) {
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
    if (game.teams.length === 3) {
      // nothing!
      return null;
    } else if (game.teams.length === 4) {
      let rightID = game.teams
        .filter(t => t.players.some(p => p.id !== myID))[0]
        .players.filter(p => p.id !== myID)[1].id;
      return game.hands[rightID];
    }
  } else if (position !== 'left') {
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
          <Grid item>{'TODO: Do action here'}</Grid>
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

  // const refreshButton = (
  //   <button
  //     onClick={() => onRefreshCurrentGame(activeGame.id)}
  //     className='hover:text-white'
  //   >
  //     Refresh
  //   </button>
  // );

  // if (!activeGame) {
  //   return (
  //     <div>
  //       This will be a page for a game, but we don't know what the game is.
  //       <br></br>
  //       {refreshButton}
  //     </div>
  //   );
  // }

  // const myColor = activeGame.player_colors[currentUser.id];
  // let gameResp = [];
  // let scoreChildren = [];
  // let playerNamesByID = {};
  // activeGame.players.forEach(player => {
  //   playerNamesByID[player.id] = player.name;
  // });

  // for (const [key, val] of Object.entries(activeGame)) {
  //   gameResp.push(`${key}: ${val} `);
  //   gameResp.push(<br key={`br ${key}`}></br>);

  //   switch (key) {
  //     case 'current_scores':
  //       let lagScores = activeGame.lag_scores;
  //       let playerColors = activeGame.player_colors;
  //       let teams = {};
  //       for (const [playerName, color] of Object.entries(playerColors)) {
  //         if (teams[color]) {
  //           teams[color] += `, ${playerName}`;
  //         } else {
  //           teams[color] = `${playerName}`;
  //         }
  //       }
  //       for (const [color, curscore] of Object.entries(val)) {
  //         let team = ` (${teams[color]})`;
  //         let scoreStr = `${color}${team}: ${curscore}`;
  //         if (lagScores && lagScores[color]) {
  //           scoreStr += ` (from ${lagScores[color]})`;
  //         }
  //         if (color === myColor) {
  //           scoreChildren.unshift(
  //             <strong key='myTeamScore'>{scoreStr}</strong>,
  //             <br key={`br ${color}`}></br>,
  //           );
  //         } else {
  //           scoreChildren.push(scoreStr, <br key={`br ${color}`}></br>);
  //         }
  //       }
  //       break;
  //   }
  // }
  // // TODO I know this component is super preliminary, but we would probably do well to break it down into more manageable subcomponents :#
  // return (
  //   <div>
  //     <div>Players are: {activeGame.players.map(p => p.name).join(', ')}</div>
  //     <div>
  //       <h2>Scores:</h2>
  //       {scoreChildren}
  //     </div>
  //     <div>Dealer: {activeGame.current_dealer}</div>
  //     <div>Phase: {activeGame.phase}</div>
  //     {!['Deal', 'BuildCrib'].includes(activeGame.phase) && (
  //       <div>Cut Card: {jsonCardToCard(activeGame.cut_card)}</div>
  //     )}
  //     {activeGame.hands[currentUser.id] ? (
  //       <div>
  //         My Hand:{' '}
  //         {activeGame.hands[currentUser.id].map(c => jsonCardToCard(c))}
  //       </div>
  //     ) : null}
  //     {Object.keys(activeGame.hands)
  //       .filter(k => k !== currentUser.id)
  //       .map(k => (
  //         <div key={k}>
  //           {k}'s Hand:{' '}
  //           {activeGame.hands[k]
  //             ? activeGame.hands[k].map(c => jsonCardToCard(c))
  //             : 'empty/unknown'}
  //         </div>
  //       ))}
  //     {activeGame.crib ? (
  //       <div>
  //         Crib: {activeGame.crib.map(c => jsonCardToCard(c).join(', '))}
  //       </div>
  //     ) : (
  //       <div>no crib yet</div>
  //     )}
  //     <br></br>
  //     This will be a page for the game of a user.
  //     <br></br>
  //     {gameResp}
  //     <br></br>
  //     {refreshButton}
  //   </div>
  // );
};

const jsonCardToCard = card => {
  // ${card.name}
  // return `${card.value} ${card.suit}`;
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
