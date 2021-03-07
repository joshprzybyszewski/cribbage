import { User } from '../../../auth/slice';

export interface ActionInputProps {
    isBlocking: boolean;
}

export interface CreateGameResponse {
    id: number;
    players: User[];
    player_colors: {
        [key: string]: string;
    };
    blocking_players: {
        [key: string]: string;
    };
    current_dealer: string;
}
