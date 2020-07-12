import React from 'react';
import { Link } from 'react-router-dom';
import * as Icons from '../icons';

export default () => {
  return (
    <div className='flex flex-col py-4 items-center'>
      <Link to='/newgame'>
        <Icons.CirclePlusIcon className='mb-4 w-8 text-gray-300 hover:text-white' />
      </Link>
      <Link to='/home'>
        <Icons.ListIcon className='mb-4 w-8 text-gray-300 hover:text-white' />
      </Link>
      <Link to='/account'>
        <Icons.UserIcon className='w-8 text-gray-300 hover:text-white' />
      </Link>
    </div>
  );
};
