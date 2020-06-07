import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { authActions } from '../../sagas/actions';

const Landing = ({ register }) => {
  const [formData, setFormData] = useState({ id: '', name: '' });

  return (
    <div className='max-w-sm m-auto mt-4'>
      <h1 className='text-4xl'>Welcome to Cribbage!</h1>
      <p>Play cribbage against your friends online. Get started now!</p>
      <form
        onSubmit={e => {
          e.preventDefault();
          register(formData);
        }}
      >
        <input
          onChange={e => setFormData({ ...formData, id: e.target.value })}
          value={formData.id}
          placeholder='Username'
          required
          className='mt-2 pl-2 h-8 shadow-sm rounded-lg block w-full focus:outline-none focus:shadow-md'
        ></input>
        <input
          onChange={e => setFormData({ ...formData, name: e.target.value })}
          value={formData.name}
          placeholder='Display name'
          required
          className='mt-2 pl-2 h-8 shadow-sm rounded-lg block w-full focus:outline-none focus:shadow-md'
        ></input>
        <p className='mt-1 text-xs text-gray-600'>
          Already have an account?{' '}
          <span>
            <Link to='/login' className='hover:text-gray-500 hover:underline'>
              Log in here.
            </Link>
          </span>
        </p>
        <input
          type='submit'
          value='register'
          className='mt-1 btn btn-primary'
        ></input>
      </form>
    </div>
  );
};

Landing.propTypes = {
  register: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    register: formData =>
      dispatch(authActions.register(formData.id, formData.name)),
  };
};

export default connect(null, mapDispatchToProps)(Landing);
