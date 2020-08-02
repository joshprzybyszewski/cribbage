import { useState } from 'react';

export const useFormInput = initialState => {
  const [val, setVal] = useState(initialState);
  const handleChange = e => setVal(e.target.value);
  return [val, handleChange];
};
