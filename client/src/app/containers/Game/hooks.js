import { selectCurrentGameID } from 'app/containers/Game/selectors';
import { selectCurrentUser } from 'auth/selectors';
import { useSelector } from 'react-redux';

export const useCurrentPlayerAndGame = () => {
  const currentUser = useSelector(selectCurrentUser);
  const gameID = useSelector(selectCurrentGameID);
  return { currentUser, gameID };
};
