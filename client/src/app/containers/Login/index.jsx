import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import { makeStyles } from '@material-ui/core/styles';
import {
  Button,
  Container,
  Link,
  TextField,
  CssBaseline,
  Typography,
} from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  title: {
    fontSize: '2rem',
  },
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

const LoginForm = () => {
  // hooks
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const [playerID, setPlayerID] = useState('');

  // event handlers
  const onSubmitLoginForm = event => {
    event.preventDefault();
    dispatch(actions.login(playerID, history));
  };
  const onInputChange = event => {
    setPlayerID(event.target.value);
  };

  const classes = useStyles();

  return (
    <Container component='main' maxWidth='sm'>
      <CssBaseline />
      <div className={classes.paper}>
        <Typography component='h1' className={classes.title}>
          Welcome to Cribbage!
        </Typography>
        <p>Play cribbage against your friends online. Get started now!</p>
        <form className={classes.form} onSubmit={onSubmitLoginForm}>
          <TextField
            variant='outlined'
            margin='normal'
            required
            fullWidth
            label='Username'
            name='id'
            autoFocus
            onChange={onInputChange}
          />
          <Button
            type='submit'
            fullWidth
            variant='contained'
            color='primary'
            className={classes.submit}
          >
            Login
          </Button>
        </form>
        <Link href='/register' variant='body2'>
          Don't have an account? Register here
        </Link>
      </div>
    </Container>
    // <div className='max-w-sm m-auto mt-12'>
    //   <h1 className='text-4xl'>Login to Cribbage</h1>
    //   <form onSubmit={onSubmitLoginForm}>
    //     <input
    //       placeholder='Username'
    //       onChange={onInputChange}
    //       className='form-input'
    //     ></input>
    //     <p className='mt-1 text-xs text-gray-600'>
    //       Don't have an account?{' '}
    //       <span>
    //         <Link
    //           to='/register'
    //           className='hover:text-gray-500 hover:underline'
    //         >
    //           Register here.
    //         </Link>
    //       </span>
    //     </p>
    //     <input
    //       type='submit'
    //       value='login'
    //       className='mt-1 btn btn-primary'
    //     ></input>
    //   </form>
    // </div>
  );
};

export default LoginForm;
