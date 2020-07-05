import React from 'react';
import { useSelector } from 'react-redux';
import { selectCurrentUser } from '../../../auth/selectors';
import ActiveGamesTable from './ActiveGamesTable';

const Home = () => {
  const currentUser = useSelector(selectCurrentUser);

  return <div>
    Welcome, {currentUser.name}!
    <ActiveGamesTable/>
  </div>;
};

export default Home;
