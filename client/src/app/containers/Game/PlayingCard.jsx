import React from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';

const PlayingCard = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const dispatch = useDispatch();

  if (!props.name) {
    return null;
  } else if (props.name === 'unknown') {
    return (
      <div
        class={`w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800`}
      ></div>
    );
  }

  // todo figure out how to state
  const chosenCards = [];

  const useRed = !['Spade', 'Clubs'].includes(props.suit);
  return (
    <div
      onClick={
        props.mine
          ? () => {
              if (!props.disabled) {
                console.log(`clicked my card ${props.name}`);
                dispatch(actions.selectCard({ card: props.card }));
              }
            }
          : () => {
              console.log(`clicked opponents card ${props.name}`);
            }
      }
      class={`w-12 h-16 text-center align-middle inline-block border-2 ${
        chosenCards.includes(props.name) ? 'border-red-700' : 'border-black'
      } ${props.disabled ? 'bg-gray-500' : 'bg-white'} ${
        useRed ? 'text-red-700' : 'text-black'
      }`}
    >
      {props.name}
    </div>
  );
};

export default PlayingCard;
