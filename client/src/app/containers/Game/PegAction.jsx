import React from 'react';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import SendIcon from '@material-ui/icons/Send';
import { gameSaga } from 'app/containers/Game/saga';
import { selectSelectedCards } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { useCurrentPlayerAndGame } from './hooks';

const PegAction = ({ isBlocking }) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const selectedCards = useSelector(selectSelectedCards);
  const { currentUser, gameID } = useCurrentPlayerAndGame();

  return (
    <ButtonGroup
      orientation='vertical'
      color='primary'
      aria-label='vertical outlined primary button group'
    >
      <Button
        disabled={!isBlocking}
        color='secondary'
        onClick={() => {
          dispatch(actions.sayGo({ userID: currentUser.id, gameID }));
        }}
      >
        Say Go
      </Button>
      <Button
        disabled={!isBlocking || selectedCards.length !== 1}
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          dispatch(
            actions.pegCard({
              userID: currentUser.id,
              gameID,
              card: selectedCards[0],
            }),
          );
        }}
      >
        Peg
      </Button>
    </ButtonGroup>
  );
};

PegAction.propTypes = {
  isBlocking: PropTypes.bool.isRequired,
};

export default PegAction;
