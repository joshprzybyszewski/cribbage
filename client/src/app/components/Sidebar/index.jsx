import React from 'react';
import { Link } from 'react-router-dom';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import HomeIcon from '@material-ui/icons/Home';
import PersonIcon from '@material-ui/icons/Person';

export default () => {
  return (
    <div className='w-12 bg-gray-400'>
      <div className='flex flex-col py-4 items-center'>
        <Link to='/home'>
          <HomeIcon
            fontSize='large'
            className='mb-4 text-blue-800 hover:text-blue-600'
          />
        </Link>
        <Link to='/newgame'>
          <AddCircleOutlineIcon
            fontSize='large'
            className='mb-4 text-blue-800 hover:text-blue-600'
          />
        </Link>
        <Link to='/account'>
          <PersonIcon
            fontSize='large'
            className='text-blue-800 hover:text-blue-600'
          />
        </Link>
      </div>
    </div>
  );
};
