export const alertTypes = {
  error: 'error',
  warning: 'warning',
  success: 'success',
  info: 'info',
};

export const alertTypeToStyle = t => {
  switch (t) {
    case alertTypes.success:
      return 'alert-success';
    case alertTypes.error:
      return 'alert-error';
    case alertTypes.warning:
      return 'alert-warning';
    default:
      return 'alert-info';
  }
};
