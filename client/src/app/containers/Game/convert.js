const overcomesMap = {
  deal: 0,
  crib: 1,
  cut: 2,
  peg: 3,
  counthand: 4,
  countcrib: 5,
};

const newPlayerAction = (myID, gameID, phase, action) => {
  return {
    pID: myID,
    gID: gameID,
    o: overcomesMap[phase],
    a: action,
  };
};

export { newPlayerAction };
