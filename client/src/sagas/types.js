// "reducer" actions should never be used outside of the reducers and the sagas pushing to the reducers
export const auth = {
  reducer: {
    LOGIN: 'LOGIN_REDUCER',
    LOGOUT: 'LOGOUT_REDUCER',
    REGISTER: 'REGISTER_REDUCER',
    LOGIN_FAILED: 'LOGIN_FAILED_REDUCER',
    REGISTER_FAILED: 'REGISTER_FAILED_REDUCER',
  },
  LOGIN: 'LOGIN',
  LOGOUT: 'LOGOUT',
  REGISTER: 'REGISTER',
};

export const alert = {
  reducer: {
    ADD_ALERT: 'ADD_ALERT_REDUCER',
    REMOVE_ALERT: 'REMOVE_ALERT_REDUCER',
  },
  ADD_ALERT: 'ADD_ALERT',
};
