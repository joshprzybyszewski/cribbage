import React from 'react';
import { Link } from 'react-router-dom';

export default () => {
  return (
    <div className='flex flex-col py-4 items-center'>
      <Link to='/newgame'>
        <CirclePlusIcon className='mb-4 text-gray-300 hover:text-white w-8' />
      </Link>
      <Link to='/home'>
        <ListIcon className='mb-4 text-gray-300 hover:text-white w-8' />
      </Link>
      <Link to='/account'>
        <UserIcon className='text-gray-300 hover:text-white w-8' />
      </Link>
    </div>
  );
};

const CirclePlusIcon = props => (
  <div {...props}>
    <svg
      xmlns='http://www.w3.org/2000/svg'
      className='icon icon-tabler icon-tabler-circle-plus stroke-current'
      viewBox='0 0 24 24'
      stroke-width='1.5'
      fill='none'
      stroke-linecap='round'
      stroke-linejoin='round'
    >
      <path stroke='none' d='M0 0h24v24H0z' />
      <circle cx='12' cy='12' r='9' />
      <line x1='9' y1='12' x2='15' y2='12' />
      <line x1='12' y1='9' x2='12' y2='15' />
    </svg>
  </div>
);

const ListIcon = props => (
  <div {...props}>
    <svg
      xmlns='http://www.w3.org/2000/svg'
      class='icon icon-tabler icon-tabler-notes stroke-current'
      viewBox='0 0 24 24'
      stroke-width='1.5'
      fill='none'
      stroke-linecap='round'
      stroke-linejoin='round'
    >
      <path stroke='none' d='M0 0h24v24H0z' />
      <rect x='5' y='3' width='14' height='18' rx='2' />
      <line x1='9' y1='7' x2='15' y2='7' />
      <line x1='9' y1='11' x2='15' y2='11' />
      <line x1='9' y1='15' x2='13' y2='15' />
    </svg>
  </div>
);

const UserIcon = props => (
  <div {...props}>
    <svg
      xmlns='http://www.w3.org/2000/svg'
      class='icon icon-tabler icon-tabler-user stroke-current'
      viewBox='0 0 24 24'
      stroke-width='1.5'
      fill='none'
      stroke-linecap='round'
      stroke-linejoin='round'
    >
      <path stroke='none' d='M0 0h24v24H0z' />
      <circle cx='12' cy='7' r='4' />
      <path d='M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2' />
    </svg>
  </div>
);
