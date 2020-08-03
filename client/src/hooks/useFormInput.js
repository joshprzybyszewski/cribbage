import { useState } from 'react';

export const useFormInput = initialState => {
  const [value, setValue] = useState(initialState);
  const handleChange = e => setValue(e.target.value);
  return [{ value, setValue }, handleChange];
};
