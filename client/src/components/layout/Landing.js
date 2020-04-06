import React, { Fragment, useState } from 'react';

export default () => {
  const [formData, setFormData] = useState({
    username: '',
    displayName: '',
  });

  const onChange = (e) =>
    setFormData({ ...formData, [e.target.name]: e.target.value });

  const onSubmit = (e) => {};

  const handleRegister = (e) => {};

  return (
    <Fragment>
      <h1 className='large text-primary'>Welcome to Cribbage!</h1>
      <form className='form' onSubmit={(e) => onSubmit(e)}>
        <div className='form-group'>
          <input
            onChange={(e) => onChange(e)}
            type='text'
            placeholder='Username'
            name='username'
            required
          />
        </div>
        <div className='form-group'>
          <input
            onChange={(e) => onChange(e)}
            type='text'
            placeholder='Display Name'
            name='displayName'
            required
          />
        </div>

        <button type='submit' className='btn btn-primary my-1'>
          Login
        </button>
        <button onClick={(e) => handleRegister(e)} className='btn my-1'>
          Register
        </button>
      </form>
    </Fragment>
  );
};
