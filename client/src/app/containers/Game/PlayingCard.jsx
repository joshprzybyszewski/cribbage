import React from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { gameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';

const PlayingCard = props => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: gameSaga });
  const [chosen, setChosen] = React.useState(false);

  if (!props.card) {
    return null;
  } else if (props.card.name === 'unknown') {
    // Currently, this returns a grayed out box, but it should show
    // a back of a card
    return (
      <div
        class={`w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800`}
      />
    );
  }

  const useRed = !['Spade', 'Clubs'].includes(props.suit);
  return (
    <div
      onClick={
        props.mine
          ? () => {
              if (!props.disabled) {
                setChosen(!chosen);
              }
            }
          : () => {}
      }
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
