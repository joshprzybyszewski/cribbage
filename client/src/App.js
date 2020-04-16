import React from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import store from './store';
import Alert from './components/layout/Alert';
import Home from './components/home/Home';
import Landing from './components/landing/Landing';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Alert />
        <section>
          <Switch>
            <Route exact path='/' component={Landing} />
            <Route exact path='/home' component={Home} />
          </Switch>
        </section>
      </Router>
    </Provider>
  );
}

export default App;
