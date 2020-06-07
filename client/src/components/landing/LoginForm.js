import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { authActions } from '../../sagas/actions';

const LoginForm = ({ login }) => {
  const [playerID, setPlayerID] = useState('');
  return (
    <div className='max-w-sm m-auto mt-4'>
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
          className='pl-2 h-8 shadow-sm rounded-lg block w-full focus:outline-none focus:shadow-md'
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
          className='mt-1 px-4 py-1 rounded-lg text-lg text-gray-300 uppercase bg-blue-800 hover:bg-blue-700 hover:text-white'
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
