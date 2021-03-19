/* eslint-disable jsx-a11y/no-static-element-interactions */
/* eslint-disable jsx-a11y/click-events-have-key-events */
// TODO don't disable eslint. maybe use a button instead
import React from 'react';

import { Card } from './models';
import { useGame } from './useGame';

interface Props {
    card: Card;
    disabled: boolean;
    mine?: boolean;
}

const PlayingCard: React.FunctionComponent<Props> = ({
    card,
    disabled,
    mine,
}) => {
    const { selectedCards, toggleSelectedCard } = useGame();
    const useRed = !['Spades', 'Clubs'].includes(card.suit);

    if (!card) {
        return null;
    }
    if (card.name === 'unknown') {
        // Currently, this returns a grayed out box, but it should show
        // a back of a card
        return (
            <div className='w-12 h-16 text-center align-middle inline-block border-2 bg-gray-800' />
        );
    }

    const chosen = selectedCards.indexOf(card) !== -1;
    const handleClick = () => {
        if (!disabled && mine) {
            toggleSelectedCard(card);
        }
    };

    return (
        <div
            onClick={handleClick}
            className={`w-12 h-16 text-center align-middle inline-block border-2 border-black ${
                disabled ? 'bg-gray-500' : 'bg-white'
            } ${useRed ? 'text-red-700' : 'text-black'}`}
            style={{
                position: 'relative',
                top: chosen ? '-10px' : '',
            }}
        >
            {card.name}
        </div>
    );
};

PlayingCard.defaultProps = {
    mine: false,
};

export default PlayingCard;
