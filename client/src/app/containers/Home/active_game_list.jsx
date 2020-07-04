import React from 'react';
import { useSelector } from 'react-redux';
import { selectCurrentUser } from '../../../auth/selectors';

const ActiveGames = () => {
  const currentUser = useSelector(selectCurrentUser);

  return <div>Your ({currentUser.name}) Active Games are: nothing! ha!</div>;
};

export default ActiveGames;
