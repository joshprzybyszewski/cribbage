import React from 'react';

const ScoreBoard = ({ teams }) => {
  return (
    <div className='py-2 px-4 bg-gray-400 border-gray-700 rounded-lg border-2'>
      Scores:{' '}
      {teams.map(t => (
        <div className='flex flex-row'>
          <div className='flex-1 capitalize'>{t.color}</div>
          <div className='flex-2'>{t.current_score}</div>
          <div className='flex-2'>
            ({t.players.map(p => p.name).join(', ')})
          </div>
        </div>
      ))}
    </div>
  );
};

export default ScoreBoard;
