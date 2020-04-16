import React, { useState } from 'react';
import { Row, Card, Col } from 'antd';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';

const Landing = props => {
  const tabList = [
    { key: 'login', tab: 'Login' },
    { key: 'register', tab: 'Register' },
  ];

  const [tabKey, setTabKey] = useState('login');
  return (
    <div className='landing'>
      <div className='dark-overlay'>
        <Row style={{ height: '33%' }}>
          <Col span={24} />
        </Row>
        <Row style={{ height: '33%' }}>
          <Col span={7} />
          <Col span={10}>
            <h1 style={{ color: '#fff' }}>Welcome to Cribbage!</h1>
            <Card
              title='Login or register to play cribbage against your friends online'
              tabList={tabList}
              activeTabKey={tabKey}
              onTabChange={k => setTabKey(k)}
            >
              {tabKey === 'login' ? <LoginForm /> : <RegisterForm />}
            </Card>
          </Col>
          <Col span={7} />
        </Row>
        <Row style={{ height: '33%' }}>
          <Col span={24} />
        </Row>
      </div>
    </div>
  );
};

export default Landing;
