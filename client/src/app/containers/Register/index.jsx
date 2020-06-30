import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { actions } from '../../../auth/slice';

const Landing = () => {
  const dispatch = useDispatch();
  const [formData, setFormData] = useState({ id: '', name: '' });
  const onSubmitForm = event => {
    event.preventDefault();
    dispatch(actions.register(formData));
  };
  const onInputChange = event =>
    setFormData({ ...formData, id: event.target.value });

  return (
    <div className='max-w-sm m-auto mt-12'>
      <h1 className='text-4xl'>Welcome to Cribbage!</h1>
      <p>Play cribbage against your friends online. Get started now!</p>
      <form onSubmit={onSubmitForm}>
        <input
          onChange={onInputChange}
          value={formData.id}
          placeholder='Username'
          required
          className='mt-2 form-input'
        ></input>
        <input
          onChange={onInputChange}
          value={formData.name}
          placeholder='Display name'
          required
          className='mt-2 form-input'
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

export default Landing;
