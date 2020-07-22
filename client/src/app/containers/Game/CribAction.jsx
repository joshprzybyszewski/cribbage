import React from 'react';

import Button from '@material-ui/core/Button';
import SendIcon from '@material-ui/icons/Send';
import { gameSaga } from 'app/containers/Game/saga';
import {
  selectCurrentAction,
  selectCurrentGame,
} from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const expNumCardsToCribForGame = game => {
  if (game.teams.length === 3 || game.teams[0].players.length === 2) {
    return 1;
  }

  return 2;
};

const CribAction = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);
  const activeGame = useSelector(selectCurrentGame);

  return (
    <Button
      disabled={
        !props.isBlocking ||
        currentAction.selectedCards.length !==
          expNumCardsToCribForGame(activeGame)
      }
      variant='contained'
      color='primary'
      endIcon={<SendIcon />}
      onClick={() => {
        dispatch(actions.buildCrib());
      }}
    >
      Build Crib
    </Button>
  );
};

export default CribAction;
