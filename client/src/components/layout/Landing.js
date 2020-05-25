import React, { useState } from 'react';
import { Card } from 'antd';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import Alert from './Alert';

const Landing = props => {
  const tabList = [
    { key: 'login', tab: 'Login' },
    { key: 'register', tab: 'Register' },
  ];

  const [tabKey, setTabKey] = useState('login');
  return (
    <section className='landing'>
      <div className='dark-overlay'>
        <div className='landing-inner'>
          <h1>Welcome to Cribbage!</h1>
          <Card
            title='Login or register to play cribbage against your friends online'
            tabList={tabList}
            activeTabKey={tabKey}
            onTabChange={k => setTabKey(k)}
          >
            <Alert />
            {tabKey === 'login' ? <LoginForm /> : <RegisterForm />}
          </Card>
        </div>
      </div>
    </section>
  );
};

export default Landing;
