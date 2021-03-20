import { User } from '../../../auth/slice';

export interface Player extends User {
    color: string;
}

export type Suit = 'Spades' | 'Clubs' | 'Diamonds' | 'Hearts';
type SuitLetter = 'C' | 'D' | 'H' | 'S';
type ValueLetter =
    | 'A'
    | '2'
    | '3'
    | '4'
    | '5'
    | '6'
    | '7'
    | '8'
    | '9'
    | '10'
    | 'J'
    | 'Q'
    | 'K';

export type Value = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13;

// prettier doesn't know about this typescript syntax so we have to disable it here
// eslint-disable-next-line prettier/prettier
export type CardName = `${ValueLetter}${SuitLetter}`

export interface Card {
    name: CardName | 'unknown';
    suit: Suit;
    value: Value;
}

function getValueLetter(val: Value): ValueLetter {
    switch (val) {
        case 1:
            return 'A';
        case 2:
            return '2';
        case 3:
            return '3';
        case 4:
            return '4';
        case 5:
            return '5';
        case 6:
            return '6';
        case 7:
            return '7';
        case 8:
            return '8';
        case 9:
            return '9';
        case 10:
            return '10';
        case 11:
            return 'J';
        case 12:
            return 'Q';
        case 13:
            return 'K';
        default:
            return 'A';
    }
}

function getCardName(val: Value, s: Suit): CardName {
    return <CardName>`${getValueLetter(val)}${s[0]}`
}

export function getCard(val: Value, s: Suit): Card {
    return {
        name: getCardName(val, s),
        suit: s,
        value: val,
    }
}

export interface PeggedCard {
    card: Card;
    player: string;
}

export interface Team {
    players: Player[];
    color: string;
    current_score: number;
    lag_score: number;
}

export interface Game {
    id: number;
    teams: Team[];
    phase: Phase;
    current_peg: number;
    blocking_players: {
        [key: string]: Blocker;
    };
    current_dealer: string;
    hands: {
        [key: string]: Card[];
    };
    crib: Card[];
    cut_card: Card;
    pegged_cards: PeggedCard[];
}

type Blocker =
    | 'DealCards'
    | 'CribCard'
    | 'CutCard'
    | 'PegCard'
    | 'CountHand'
    | 'CountCrib'
    | 'unknownBlocker';

export type Phase =
    | 'Deal'
    | 'BuildCrib'
    | 'Cut'
    | 'Pegging'
    | 'Counting'
    | 'CribCounting'
    | 'unknownPhase';
