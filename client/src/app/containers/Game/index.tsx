import React from 'react';

import { Grid, IconButton } from '@material-ui/core';
import RefreshIcon from '@material-ui/icons/Refresh';

import { useAuth } from '../../../auth/useAuth';
import ActionBox from './ActionBox';
import CribHand from './CribHand';
import { Game, Phase } from './models';
import PlayerHand from './PlayerHand';
import PlayingCard from './PlayingCard';
import ScoreBoard from './ScoreBoard';
import { useGame } from './useGame';

const showCutCard = (phase: Phase) =>
    ['Pegging', 'CountHands', 'CountCrib'].includes(phase);

const hasDealtHands = (phase: Phase) =>
    !['unknownPhase', 'Deal', 'DealingReady'].includes(phase);

const handForPlayer = (
    game: Game,
    myID: string,
    position: 'across' | 'right' | 'left',
) => {
    if (!hasDealtHands(game.phase)) {
        return [];
    }
    const numPlayers = game.teams.reduce(
        (prev, team) => prev + team.players.length,
        0,
    );
    const isFourPlayer = numPlayers === 4;
    if (position === 'across') {
        if (game.teams.length === 3) {
            const secondPlayerID = game.teams.filter(
                t => !t.players.some(p => p.id === myID),
            )[1].players[0].id;
            return game.hands[secondPlayerID] ?? [];
        }
        if (isFourPlayer) {
            const partnerID = game.teams
                .filter(t => t.players.some(p => p.id === myID))[0]
                .players.filter(p => p.id !== myID)[0].id;
            return game.hands[partnerID] ?? [];
        }
        const opponentID = game.teams.filter(
            t => !t.players.some(p => p.id === myID),
        )[0].players[0].id;
        return game.hands[opponentID] ?? [];
    }
    if (position === 'right') {
        if (isFourPlayer) {
            const rightID = game.teams
                .filter(t => t.players.some(p => p.id !== myID))[0]
                .players.filter(p => p.id !== myID)[1].id;
            return game.hands[rightID] ?? [];
        }
        // nothing!
        return [];
    }
    if (position !== 'left' || !isFourPlayer) {
        return [];
    }
    // position is left
    const leftID = game.teams
        .filter(t => t.players.some(p => p.id !== myID))[0]
        .players.filter(p => p.id !== myID)[0].id;
    return game.hands[leftID] ?? [];
};

const GamePage: React.FunctionComponent = () => {
    const { game, refreshGame } = useGame();
    const { currentUser } = useAuth();
    const myHand = hasDealtHands(game.phase) ? game.hands[currentUser.id] : [];

    return (
        <Grid container xl spacing={1} direction='row' justify='space-between'>
            <Grid
                item
                container
                md
                spacing={2}
                direction='column'
                align-content='space-between'
            >
                <Grid item xs sm container>
                    <PlayerHand
                        hand={handForPlayer(game, currentUser.id, 'across')}
                    />
                </Grid>
                <Grid
                    item
                    xs
                    md
                    container
                    justify='space-between'
                    align-content='center'
                >
                    <Grid item>
                        <PlayerHand
                            side
                            hand={handForPlayer(game, currentUser.id, 'left')}
                        />
                    </Grid>
                    <Grid item>
                        <ActionBox
                            phase={game.phase}
                            isBlocking={Object.keys(
                                game.blocking_players,
                            ).includes(currentUser.id)}
                        />
                    </Grid>
                    <Grid item>
                        <PlayerHand
                            side
                            hand={handForPlayer(game, currentUser.id, 'right')}
                        />
                    </Grid>
                </Grid>
                <Grid item xs sm container>
                    <PlayerHand mine hand={myHand} />
                </Grid>
            </Grid>
            <Grid item container xs direction='column' spacing={1}>
                <Grid item>
                    <IconButton aria-label='refresh' onClick={refreshGame}>
                        <RefreshIcon />
                    </IconButton>
                    <ScoreBoard
                        teams={game.teams}
                        current_dealer={game.current_dealer}
                    />
                </Grid>
                <Grid item>
                    {[
                        showCutCard(game.phase) ? (
                            <PlayingCard
                                key='cutCard'
                                card={game.cut_card}
                                disabled={false}
                                mine={false}
                            />
                        ) : (
                            <div key='deckTODOdiv'>
                                TODO put an image of the deck here
                            </div>
                        ),
                        <CribHand key='cribHand' cards={game.crib} />,
                        <div key='currentPeg'>
                            {game.phase === 'Pegging'
                                ? `Current Peg: ${
                                      game.current_peg ? game.current_peg : 0
                                  }`
                                : ''}
                        </div>,
                    ]}
                </Grid>
            </Grid>
        </Grid>
    );
};

export default GamePage;
