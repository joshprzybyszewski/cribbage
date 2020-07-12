import React from 'react';
import { Link } from 'react-router-dom';
import * as Icons from '../icons';

export default () => {
  return (
    <div className='w-12 bg-gray-400'>
      <div className='flex flex-col py-4 items-center'>
        <Link to='/newgame'>
          <Icons.CirclePlusIcon className='mb-4 w-8 text-blue-800 hover:text-blue-600' />
        </Link>
        <Link to='/home'>
          <Icons.ListIcon className='mb-4 w-8 text-blue-800 hover:text-blue-600' />
        </Link>
        <Link to='/account'>
          <Icons.UserIcon className='w-8 text-blue-800 hover:text-blue-600' />
        </Link>
      </div>
    </div>
  );
};
