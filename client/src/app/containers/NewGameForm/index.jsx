import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import makeStyles from '@material-ui/core/styles/makeStyles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

import { createGameSaga } from './saga';
import { sliceKey, reducer, actions } from './slice';

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

const NewGameForm = () => {
  // hooks
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: createGameSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const [opp1ID, setOpp1ID] = useState('');
  const [opp2ID, setOpp2ID] = useState('');
  const [teammateID, setTeammateID] = useState('');

  // event handlers
  const onSubmitLoginForm = event => {
    event.preventDefault();
    dispatch(actions.createGame(opp1ID, opp2ID, teammateID, history));
  };
  const onOpp1IDChange = event => {
    setOpp1ID(event.target.value);
  };
  const onOpp2IDChange = event => {
    setOpp2ID(event.target.value);
  };
  const onTeammateIDChange = event => {
    setTeammateID(event.target.value);
  };

  const classes = useStyles();

  return (
    <Container component='main' maxWidth='sm'>
      <CssBaseline />
      <div className={classes.paper}>
        <Typography component='h1' className={classes.title}>
          Start Game
        </Typography>
        <form className={classes.form} onSubmit={onSubmitLoginForm}>
          <TextField
            variant='outlined'
            margin='normal'
            required
            fullWidth
            label='Opponent 1'
            name='id1'
            autoFocus
            onChange={onOpp1IDChange}
          />
          <TextField
            disabled
            variant='outlined'
            margin='normal'
            fullWidth
            label='Opponent 2'
            name='id2'
            autoFocus
            onChange={onOpp2IDChange}
          />
          <TextField
            disabled
            variant='outlined'
            margin='normal'
            fullWidth
            label='Teammate'
            name='teammateID'
            autoFocus
            onChange={onTeammateIDChange}
          />
          <Button
            type='submit'
            fullWidth
            variant='contained'
            color='primary'
            className={classes.submit}
            disabled={opp1ID.length === 0}
          >
            Challenge
          </Button>
        </form>
      </div>
    </Container>
  );
};

export default NewGameForm;
