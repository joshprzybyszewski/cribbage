import React from 'react';

import { selectSuggestions } from 'app/containers/Suggestions/selectors';
import { sliceKey, reducer } from 'app/containers/Suggestions/slice';
import { useSelector } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

const Suggestions = () => {
  useInjectReducer({ key: sliceKey, reducer });

  const sugs = useSelector(selectSuggestions);

  return (
    <div className='fixed w-half-screen'>
      <div> hello there {sugs}</div>
    </div>
  );
};

export default Suggestions;
