import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { login } from '../../sagas/auth';

const LoginForm = ({ login }) => {
  return (
    <Form onFinish={formData => login(formData)}>
      <Form.Item
        name='username'
        label='Username'
        rules={[{ required: true, message: 'Please input your username!' }]}
      >
        <Input placeholder='Username' prefix={<UserOutlined />} />
      </Form.Item>
      <Form.Item>
        <Button type='primary' htmlType='submit'>
          Login
        </Button>
      </Form.Item>
    </Form>
  );
};

LoginForm.propTypes = {
  login: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    login: formData => dispatch(login(formData.username)),
  };
};

export default connect(null, mapDispatchToProps)(LoginForm);