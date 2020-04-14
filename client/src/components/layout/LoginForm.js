import React, { Fragment, useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';

const LoginForm = props => {
  return (
    <Form onFinish={vals => console.log(vals)}>
      <Form.Item name="username" label="Username">
        <Input placeholder="Username" prefix={<UserOutlined />} />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Login
        </Button>
      </Form.Item>
    </Form>
  );
};

LoginForm.propTypes = {};

export default LoginForm;
