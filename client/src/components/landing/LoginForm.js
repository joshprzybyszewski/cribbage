import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { loginAction } from '../../sagas/auth';

const LoginForm = ({ loginAsync }) => {
  return (
    <Form onFinish={formData => loginAsync(formData)}>
      <Form.Item
        name='id'
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
  loginAsync: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    loginAsync: formData => dispatch(loginAction(formData.id)),
  };
};

export default connect(null, mapDispatchToProps)(LoginForm);
