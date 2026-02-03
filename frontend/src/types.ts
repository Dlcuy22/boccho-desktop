/*
types.ts - TypeScript interfaces for Wails bindings

Interfaces:
- CharacterInfo: Character metadata from Go backend
- CharacterWindowInfo: Active window information with scale
*/

export interface CharacterInfo {
  name: string;
  path: string;
  previewPath: string;
  frameCount: number;
}

export interface CharacterWindowInfo {
  id: string;
  characterName: string;
  isRunning: boolean;
  scale: number;
}
