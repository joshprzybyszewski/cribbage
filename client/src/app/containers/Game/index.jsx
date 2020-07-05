import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { selectCurrentUser } from '../../../auth/selectors';
import { gameSaga } from '../../../game/saga';
import { sliceKey, reducer, actions } from '../../../game/slice';
import { selectCurrentGame } from '../../../game/selectors';

const Game = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const activeGame = useSelector(selectCurrentGame);
  const activeGameID = activeGame.id;
  
  // event handlers
  const onRefreshCurrentGame = () => {
    dispatch(actions.refreshGame(activeGameID, history));
  };
  const refreshButton = <button
    onClick={onRefreshCurrentGame}
    className='hover:text-white'
    >
    Refresh
  </button>;

  if ( !activeGame ) {
    return <div>
      This will be a page for a game, but we don't know what the game is.
      <br></br>
      {refreshButton}
    </div>;
  }

  let gameResp = [];
  let gameDesc = 'Players are: ';
  let scoreChildren = [];
  let dealerDesc = 'Dealer: ';
  let phaseDesc = 'Phase: ';
  let cutCardDiv;
  for (const [key, val] of Object.entries(activeGame)) {
    gameResp.push(`${key}: ${val} `);
    gameResp.push(<br key={`br ${key}`}></br>);

    switch (key) {
        case 'players':
            val.forEach((player, index) => {
                gameDesc += player.name;
                if ( index < val.length - 1 ) {
                    gameDesc += ', ';
                }
            });
            break;
        case 'current_scores':
            let lagScores = activeGame['lag_scores'];
            let playerColors = activeGame['player_colors'];
            let teams = {};
            for (const [playerName, color] of Object.entries(playerColors)) {
                if ( teams[color] ) {
                    teams[color] += `, ${playerName}`
                } else {
                    teams[color] = `${playerName}`
                }
            }
            for (const [color, curscore] of Object.entries(val)) {
                let team = ` (${teams[color]})`;
                let scoreStr = `${color}${team}: ${curscore}`;
                if (lagScores && lagScores[color]) {
                    scoreStr += ` (from ${lagScores[color]})`
                }
                scoreChildren.push(scoreStr);
                scoreChildren.push(<br key={`br ${color}`}></br>);
            }
            break;
        case 'current_dealer':
            dealerDesc += val;
            break;
        case 'phase':
            phaseDesc += val;
            break;
        case 'cut_card':
            let skip = false;
            switch (activeGame['phase']) {
                case 'Deal':
                case 'BuildCrib':
                    skip = true;
                    break;
            }
            if ( skip ) {
                break;
            }
            cutCardDiv = <div key={'cutCardDiv'}>
                Cut Card: {val.name} ({val.value} of {val.suit})
            </div>;
            break;
    }
  }

  return <div>
      <div key={'gameDescDiv'}>
        {gameDesc}
      </div>
      <div key={'scoresDiv'}>
          <h2>Scores:</h2>
        {scoreChildren}
      </div>
      <div key={'dealerDiv'}>
        {dealerDesc}
      </div>
      <div key={'phaseDiv'}>
        {phaseDesc}
      </div>
      {cutCardDiv}

      <br></br>
      This will be a page for the game of a user.
      <br></br>
      {gameResp}
      <br></br>
      {refreshButton}
    </div>;
};

export default Game;
