import React, { useState } from 'react';

import Button from '@material-ui/core/Button';
import InputAdornment from '@material-ui/core/InputAdornment';
import MenuItem from '@material-ui/core/MenuItem';
import AndroidIcon from '@material-ui/icons/Android';
import PersonIcon from '@material-ui/icons/Person';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import makeStyles from '@material-ui/core/styles/makeStyles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { createGameSaga } from 'app/containers/NewGameForm/saga';
import { sliceKey, reducer, actions } from 'app/containers/NewGameForm/slice';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';

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

  // event handlers
  const [formData, setFormData] = useState({
    id1: '',
    id2: '',
    teammateID: '',
  });
  const onFormDataChange = e => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };
  const [robots, setRobots] = useState({
    robot1: false,
    robot2: false,
    robot3: false,
  });
  const onSubmitLoginForm = event => {
    event.preventDefault();
    dispatch(
      actions.createGame(
        formData.id1,
        formData.id2,
        formData.teammateID,
        history,
      ),
    );
  };

  const robotPlayers = [
    { id: 'DumbNPC', label: 'easy' },
    { id: 'SimpleNPC', label: 'medium' },
    { id: 'CalculatedNPC', label: 'hard' },
  ];

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
            onChange={onFormDataChange}
            helperText='Please select your opponent'
            select={robots.robot1}
            InputProps={{
              endAdornment: (
                <InputAdornment position='end'>
                  {robots.robot1 ? (
                    <AndroidIcon
                      onClick={() => {
                        setRobots({ ...robots, robot1: false });
                      }}
                    />
                  ) : (
                    <PersonIcon
                      onClick={() => {
                        setRobots({ ...robots, robot1: true });
                      }}
                    />
                  )}
                </InputAdornment>
              ),
            }}
          >
            {robots.robot1 &&
              robotPlayers.map(option => (
                <MenuItem key={option.id} value={option.id}>
                  {option.label}
                </MenuItem>
              ))}
          </TextField>
          <TextField
            disabled
            variant='outlined'
            margin='normal'
            fullWidth
            label='Opponent 2'
            name='id2'
            autoFocus
            onChange={onFormDataChange}
          />
          <TextField
            disabled
            variant='outlined'
            margin='normal'
            fullWidth
            label='Teammate'
            name='teammateID'
            autoFocus
            onChange={onFormDataChange}
          />
          <Button
            type='submit'
            fullWidth
            variant='contained'
            color='primary'
            className={classes.submit}
            disabled={formData.id1.length === 0}
          >
            Challenge
          </Button>
        </form>
      </div>
    </Container>
  );
};

export default NewGameForm;
