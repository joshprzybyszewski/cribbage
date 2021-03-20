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

// TODO convert the name field to a method?

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
