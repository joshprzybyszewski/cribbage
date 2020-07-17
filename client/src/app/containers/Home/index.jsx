import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import { selectCurrentUser } from '../../../auth/selectors';
import { actions as homeActions } from './slice';

import ActiveGamesTable from './ActiveGamesTable';

const Home = () => {
  const dispatch = useDispatch();

  const currentUser = useSelector(selectCurrentUser);

  // because we pass nothing as an effect dependency (the second arg),
  // this will run once when we first render Home
  useEffect(() => {
    dispatch(homeActions.refreshActiveGames({ id: currentUser.id }));
  }, []);

  return (
    <div>
      Welcome, {currentUser.name}!
      <ActiveGamesTable />
    </div>
  );
};

export default Home;
