import React from 'react';

const PlayingCard = props => {
  if (!props.name) {
    return null;
  }

  // todo figure out how to state
  const chosenCards = [];

  const useRed = !['Spade', 'Clubs'].includes(props.suit);
  return (
    <div
      onClick={
        props.mine
          ? () => {
              // TODO select this card (if allowed?)
              console.log(`clicked my card ${props.name}`);
            }
          : () => {
              console.log(`clicked opponents card ${props.name}`);
            }
      }
      class={`w-12 h-16 text-center align-middle inline-block border-2 ${
        chosenCards.includes(props.name) ? 'border-red-700' : 'border-black'
      } bg-white ${useRed ? 'text-red-700' : 'text-black'}`}
    >
      {props.name}
    </div>
  );
};

export default PlayingCard;
