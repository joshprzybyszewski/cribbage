import React from 'react';

import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import SendIcon from '@material-ui/icons/Send';
import { gameSaga } from 'app/containers/Game/saga';
import { selectCurrentAction } from 'app/containers/Game/selectors';
import { sliceKey, reducer, actions } from 'app/containers/Game/slice';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const PegAction = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });

  const dispatch = useDispatch();

  const currentAction = useSelector(selectCurrentAction);

  return (
    <ButtonGroup
      orientation='vertical'
      color='primary'
      aria-label='vertical outlined primary button group'
    >
      <Button
        disabled={!props.isBlocking}
        color='secondary'
        onClick={() => {
          dispatch(actions.pegCard());
        }}
      >
        Say Go
      </Button>
      <Button
        disabled={!props.isBlocking || currentAction.selectedCards.length !== 1}
        color='primary'
        endIcon={<SendIcon />}
        onClick={() => {
          dispatch(actions.pegCard());
        }}
      >
        Peg
      </Button>
    </ButtonGroup>
  );
};

export default PegAction;
