import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { Link, useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';

const LoginForm = () => {
  // hooks
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const [playerID, setPlayerID] = useState('');

  // event handlers
  const onSubmitLoginForm = event => {
    event.preventDefault();
    dispatch(actions.login(playerID, history));
  };
  const onInputChange = event => {
    setPlayerID(event.target.value);
  };

  return (
    <div className='max-w-sm m-auto mt-12'>
      <h1 className='text-4xl'>Login to Cribbage</h1>
      <form onSubmit={onSubmitLoginForm}>
        <input
          placeholder='Username'
          onChange={onInputChange}
          className='form-input'
        ></input>
        <p className='mt-1 text-xs text-gray-600'>
          Don't have an account?{' '}
          <span>
            <Link
              to='/register'
              className='hover:text-gray-500 hover:underline'
            >
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
