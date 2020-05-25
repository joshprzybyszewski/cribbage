import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Alert } from 'antd';

const Alert = ({ alerts }) => {
  return (
    <div>
      {alerts.map(a => (
        <Alert key={a.id} message={a.msg} type={a.type} />
      ))}
    </div>
  );
};

Alert.propTypes = {
  alerts: PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
  alerts: state.alert,
});

export default connect(mapStateToProps, null)(Alert);
