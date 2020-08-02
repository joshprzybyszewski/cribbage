const newPlayerAction = (myID, gameID, phase, action) => {
  return {
    pID: myID,
    gID: gameID,
    o: phase,
    a: action,
  };
};

export { newPlayerAction };
