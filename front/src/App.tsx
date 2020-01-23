import React, { useRef, useEffect, useState } from 'react';
import {IonPhaser} from '@ion-phaser/react';

import './App.css';
import { PhaserApp } from './phaser/app';


const App: React.FC = () => {
  const ref = useRef<HTMLDivElement>(null)
  const [game, setGame] = useState<Phaser.Game | undefined>(undefined)
  useEffect(() => {
    if (ref.current && !game) {
      setGame(new PhaserApp(ref.current))
    }
  }, [ref, setGame, game]);
  const props: any = {
    initialize: false
  }
  return (
    <div className="App">
      <div ref={ref} >
        <IonPhaser game={game} {...props} />
      </div>
    </div>
  );
}

export default App;
