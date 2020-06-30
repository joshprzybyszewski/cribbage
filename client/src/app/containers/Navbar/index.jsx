import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Link } from 'react-router-dom';
import { selectLoggedIn } from '../../../auth/selectors';
import { actions } from '../../../auth/slice';

const Navbar = () => {
  const loggedIn = useSelector(selectLoggedIn);
  const dispatch = useDispatch();
  const onClickLogout = () => {
    dispatch(actions.logout());
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
