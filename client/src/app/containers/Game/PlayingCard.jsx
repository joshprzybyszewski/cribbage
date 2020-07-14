import React from 'react';

const PlayingCard = props => {
  if (!props.name) {
    return null;
  }
  const useRed = !['Spade', 'Clubs'].includes(props.suit);
  return (
    <div
      class={`w-12 h-16 text-center align-middle inline-block border-2 border-black bg-white ${
        useRed ? 'text-red-700' : 'text-black'
      }`}
    >
      {props.name}
    </div>
  );
};

export default PlayingCard;
