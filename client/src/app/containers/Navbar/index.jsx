import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { Link, useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { selectLoggedIn } from '../../../auth/selectors';

const Navbar = () => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const loggedIn = useSelector(selectLoggedIn);
  const history = useHistory();
  const dispatch = useDispatch();
  const onClickLogout = () => {
    dispatch(actions.logout(history));
  };

  return (
    <nav className='h-12 px-4 bg-blue-900 flex justify-between items-center text-gray-400'>
      <Link
        to={loggedIn ? '/home' : '/'}
        className='uppercase text-xl tracking-wider hover:text-white'
      >
        Cribbage
      </Link>
      {!loggedIn ? (
        <div className='flex'>
          <Link to='/login' className='px-2 hover:text-white'>
            Login
          </Link>
          <Link to='/' className='px-2 hover:text-white'>
            Register
          </Link>
        </div>
      ) : (
        <button
          onClick={onClickLogout}
          className='focus:outline-none hover:text-white'
        >
          Logout
        </button>
      )}
    </nav>
  );
};

export default Navbar;
