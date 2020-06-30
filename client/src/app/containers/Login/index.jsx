import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { Link } from 'react-router-dom';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { authSaga } from '../../../auth/saga';

const LoginForm = () => {
  // The redux stuff for auth lives outside of containers/ because it's used in a lot of places.
  // We have to inject its reducer and saga _somewhere_ and I arbitrarily chose here.
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });

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
