import React from 'react';
import { Provider } from 'react-redux';
import store from './store';
import Landing from './components/layout/Landing';
import Layout from './components/layout/Layout';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <Layout>
        <Landing />
      </Layout>
    </Provider>
  );
}

export default App;
