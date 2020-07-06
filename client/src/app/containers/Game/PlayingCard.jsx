import React from 'react';

const PlayingCard = props => {
  if (!props.name) {
    return null;
  }
  let useRed = true;
  if (props.suit === 'Spades' || props.suit === 'Clubs') {
    useRed = false;
  }
  return (
    <div
      class={`w-16 h-20 text-center align-middle inline-block border-2 border-black bg-white ${
        useRed ? 'text-red-700' : 'text-black'
      }`}
    >
      {props.name}
    </div>
  );
};

export default PlayingCard;
