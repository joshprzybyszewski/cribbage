import React, { useEffect } from 'react';

import { useAuth } from '../../../auth/useAuth';
import ActiveGamesTable from './ActiveGamesTable';
import { useActiveGames } from './useActiveGames';

const Home = () => {
    const { currentUser } = useAuth();
    const { refreshGames } = useActiveGames();

    // because we pass nothing as an effect dependency (the second arg),
    // this will run once when we first render Home
    useEffect(() => {
        refreshGames();
    }, []);

    return (
        <div>
            Welcome, {currentUser.name}!
            <ActiveGamesTable />
        </div>
    );
};

export default Home;
