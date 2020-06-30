import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { Link } from 'react-router-dom';
import { actions } from '../../../auth/slice';

const LoginForm = () => {
  const dispatch = useDispatch();
  const onLoginFormSubmit = event => {
    event.preventDefault();
    dispatch(actions.login(playerID));
  };
  const onUserIDInputChange = event => {
    setPlayerID(event.target.value);
  };

  const [playerID, setPlayerID] = useState('');
  return (
    <div className='max-w-sm m-auto mt-12'>
      <h1 className='text-4xl'>Login to Cribbage</h1>
      <form onSubmit={onLoginFormSubmit}>
        <input
          placeholder='Username'
          onChange={onUserIDInputChange}
          className='form-input'
        ></input>
        <p className='mt-1 text-xs text-gray-600'>
          Don't have an account?{' '}
          <span>
            <Link to='/' className='hover:text-gray-500 hover:underline'>
              Register here.
            </Link>
          </span>
        </p>
        <input
          type='submit'
          value='login'
          className='mt-1 btn btn-primary'
        ></input>
      </form>
    </div>
  );
};

export default LoginForm;
