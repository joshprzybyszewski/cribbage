/*
We are using tabler-icons (https://github.com/tabler/tabler-icons) under the MIT license:


MIT License

Copyright (c) 2020 Paweł Kuna

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import React from 'react';

export const CirclePlusIcon = props => (
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

export const ListIcon = props => (
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

export const UserIcon = props => (
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
