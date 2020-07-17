import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';
import { selectCurrentAction } from './selectors';

const PlayingCard = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const dispatch = useDispatch();
  const currentAction = useSelector(selectCurrentAction);

  if (!props.card) {
    return null;
  } else if (props.card.name === 'unknown') {
    // Currently, this returns a grayed out box, but it should show
    // a back of a card
    return (
      <div
        className={`w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800`}
      />
    );
  }

  let chosen = currentAction.selectedCards.indexOf(props.card) !== -1;
  let toggleChosen = () => {
    if (!props.disabled) {
      dispatch(actions.selectCard(props.card));
    }
  };

  const useRed = !['Spades', 'Clubs'].includes(props.card.suit);
  return (
    <div
      onClick={props.mine ? toggleChosen : () => {}}
      className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
        props.disabled ? 'bg-gray-500' : 'bg-white'
      } ${useRed ? 'text-red-700' : 'text-black'}`}
      style={{
        position: 'relative',
        top: chosen ? '-10px' : '',
      }}
    >
      {props.card.name}
    </div>
  );
};

export default PlayingCard;
