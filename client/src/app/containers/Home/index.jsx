import React from 'react';
import { useSelector } from 'react-redux';
import { selectCurrentUser } from '../../../auth/selectors';

const Home = () => {
  const currentUser = useSelector(selectCurrentUser);

  return <div>Welcome, {currentUser.name}!</div>;
};

export default Home;
