import { User } from '../../../auth/slice';

export interface ActionInputProps {
    isBlocking: boolean;
}

export interface CreateGameResponse {
    id: number;
    players: User[];
    // maps playerId -> peg color
    player_colors: {
        [key: string]: string;
    };
    // maps playerId -> blocking reason
    blocking_players: {
        [key: string]: string;
    };
    current_dealer: string;
}
