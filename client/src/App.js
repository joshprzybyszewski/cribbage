import React from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { Layout } from 'antd';
import store from './store';
import Alert from './components/layout/Alert';
import Home from './components/home/Home';
import Landing from './components/landing/Landing';
import './App.css';

const { Content, Header } = Layout;

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Layout>
          <Header>
            <h1 style={{ color: '#fff' }}>CRIBBAGE</h1>
          </Header>
          <Alert />
          <Switch>
            <Route exact path='/' component={Landing} />
          </Switch>
          <Content className='content'>
            <Switch>
              <Route exact path='/home' component={Home} />
            </Switch>
          </Content>
        </Layout>
      </Router>
    </Provider>
  );
}

export default App;
