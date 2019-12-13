import Phaser from 'phaser';

import {StartScene} from './scenes/Start';

var config = {
  type: Phaser.AUTO,
  width: 800,
  height: 600
};
class PhaserApp extends Phaser.Game {
  constructor() {
    super(config)
    this.scene.add('Start', StartScene);
    this.scene.start('Start')
  }
}

var game = new PhaserApp();
