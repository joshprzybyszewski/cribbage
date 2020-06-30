import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { authActions } from '../../sagas/actions';
import { actions } from '../../../auth/slice';

const LoginForm = ({ login }) => {
  const dispatch = useDispatch();
  const onLoginFormSubmit = event => {
    event.preventDefault();
    dispatch(actions.login(playerID));
  };

  const [playerID, setPlayerID] = useState('');
  return (
    <div className='max-w-sm m-auto mt-12'>
      <h1 className='text-4xl'>Login to Cribbage</h1>
      <form
        onSubmit={e => {
          e.preventDefault();
          login(playerID);
        }}
      >
        <input
          placeholder='Username'
          onChange={e => setPlayerID(e.target.value)}
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

LoginForm.propTypes = {
  login: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    login: id => dispatch(authActions.login(id)),
  };
};

export default connect(null, mapDispatchToProps)(LoginForm);
