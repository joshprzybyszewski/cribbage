import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Button, Input, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { register } from '../../sagas/auth';

const RegisterForm = ({ register, history }) => {
  const layout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 18 },
  };
  const tailLayout = {
    wrapperCol: { offset: 6, span: 18 },
  };
  return (
    <Form {...layout} onFinish={formData => register(formData, history)}>
      <Form.Item
        name='username'
        label='Username'
        rules={[{ required: true, message: 'Please input your username!' }]}
      >
        <Input placeholder='Username' prefix={<UserOutlined />} />
      </Form.Item>
      <Form.Item
        name='displayName'
        label='Display Name'
        rules={[{ required: true, message: 'Please input your display name!' }]}
      >
        <Input placeholder='Display name' />
      </Form.Item>
      <Form.Item {...tailLayout}>
        <Button type='primary' htmlType='submit'>
          Register
        </Button>
      </Form.Item>
    </Form>
  );
};

RegisterForm.propTypes = {
  register: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    register: (formData, history) =>
      dispatch(register(formData.username, formData.displayName, history)),
  };
};

export default connect(null, mapDispatchToProps)(withRouter(RegisterForm));
