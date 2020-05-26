import React from 'react';
import { Provider } from 'react-redux';
import store from './store';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import Home from './components/Home';
import Landing from './components/landing/Landing';
import Layout from './components/layout/Layout';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Layout>
          <Switch>
            <Route path='/' component={Landing} />
            <Route path='/home' component={Home} />
          </Switch>
        </Layout>
      </Router>
    </Provider>
  );
}

export default App;
