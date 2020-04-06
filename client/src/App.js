import React, { Fragment } from 'react';
import { Provider } from 'react-redux';
import store from './store';
import Landing from './components/layout/Landing';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <Fragment>
        <section className='container'>
          <Landing />
        </section>
      </Fragment>
    </Provider>
  );
}

export default App;
