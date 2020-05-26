import React, { useState } from 'react';
import { Card } from 'antd';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import './Landing.css';

const Landing = props => {
  const tabList = [
    { key: 'login', tab: 'Login' },
    { key: 'register', tab: 'Register' },
  ];

  const [tabKey, setTabKey] = useState('login');
  return (
    <div className='landing'>
      <div className='landing-content'>
        <h1>Welcome to Cribbage!</h1>
        <Card
          title='Login or register to play cribbage against your friends online'
          tabList={tabList}
          activeTabKey={tabKey}
          onTabChange={k => setTabKey(k)}
        >
          {tabKey === 'login' ? <LoginForm /> : <RegisterForm />}
        </Card>
      </div>
    </div>
  );
};

export default Landing;
