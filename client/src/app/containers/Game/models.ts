import { User } from '../../../auth/slice';

export interface Player extends User {
    color: string;
}

type Suit = 'Spades' | 'Clubs' | 'Diamonds' | 'Hearts';
type SuitLetter = 'c' | 'd' | 'h' | 's';
type ValueLetter =
    | 'a'
    | '2'
    | '3'
    | '4'
    | '5'
    | '6'
    | '7'
    | '8'
    | '9'
    | '10'
    | 'j'
    | 'q'
    | 'k';

type Value = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13;

// prettier doesn't know about this typescript syntax so we have to disable it here
// eslint-disable-next-line prettier/prettier
export type CardName = `${ValueLetter}${SuitLetter}`

export interface Card {
    name: CardName | 'unknown';
    suit: Suit;
    value: Value;
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
