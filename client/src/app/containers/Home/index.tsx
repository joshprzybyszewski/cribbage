import React, { useEffect } from 'react';

import ActiveGamesTable from 'app/containers/Home/ActiveGamesTable';
import { actions as homeActions } from 'app/containers/Home/slice';
import { selectCurrentUser } from 'auth/selectors';
import { useSelector, useDispatch } from 'react-redux';

const Home = () => {
    const dispatch = useDispatch();

    const currentUser = useSelector(selectCurrentUser);

    // because we pass nothing as an effect dependency (the second arg),
    // this will run once when we first render Home
    useEffect(() => {
        dispatch(homeActions.refreshActiveGames({ id: currentUser.id }));
    }, [dispatch, currentUser.id]);

    return (
        <div>
            Welcome, {currentUser.name}!
            <ActiveGamesTable />
        </div>
    );
};

export default Home;
