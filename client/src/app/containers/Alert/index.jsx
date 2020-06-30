import React from 'react';
import PropTypes from 'prop-types';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { connect } from 'react-redux';

const Alert = ({ alerts }) => {
  const alertTypeToStyle = t => {
    switch (t) {
      case 'success':
        return 'alert-success';
      case 'error':
        return 'alert-error';
      case 'warning':
        return 'alert-warning';
      default:
        return 'alert-info';
    }
  };
  return (
    <div className='fixed w-screen px-3 py-2'>
      {alerts.map(a => (
        // Issue#61 think about only displaying the last alert
        <div key={a.id} className={`alert ${alertTypeToStyle(a.type)} mb-2`}>
          {a.msg}
        </div>
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
