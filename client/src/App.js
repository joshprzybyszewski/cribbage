import React from 'react';
import { Provider } from 'react-redux';
import store from './store';
import Alert from './components/layout/Alert';
import Landing from './components/landing/Landing';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <Alert />
      <Landing />
    </Provider>
  );
}

export default App;
