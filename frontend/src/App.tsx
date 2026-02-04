/*
App.tsx - Main character manager UI

Components:
- App: Main application component with character grid and active windows list
- AnimatedPreview: Component that cycles through frames for animation preview
- Character cards with animated preview images
- Active windows with scale slider
*/

import { useState, useEffect, useCallback, useRef } from 'react';
import './App.css';
import { CharacterInfo, CharacterWindowInfo } from './types';
import {
  GetCharacters,
  SpawnCharacter,
  DestroyCharacter,
  GetActiveWindows,
  SetCharacterScale,
  GetPreviewFrames,
  OpenFramesDir,
  OpenConfig,
} from '../wailsjs/go/main/App';

const MAX_PREVIEW_FRAMES = 16;
const ANIMATION_FPS = 12;

function AnimatedPreview({ characterName }: { characterName: string }) {
  const [frames, setFrames] = useState<string[]>([]);
  const [currentFrame, setCurrentFrame] = useState(0);
  const [isHovered, setIsHovered] = useState(false);
  const intervalRef = useRef<number | null>(null);

  useEffect(() => {
    GetPreviewFrames(characterName, MAX_PREVIEW_FRAMES).then((f) => {
      setFrames(f || []);
    });
  }, [characterName]);

  useEffect(() => {
    if (isHovered && frames.length > 1) {
      intervalRef.current = window.setInterval(() => {
        setCurrentFrame((prev) => (prev + 1) % frames.length);
      }, 1000 / ANIMATION_FPS);
    } else {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      setCurrentFrame(0);
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isHovered, frames.length]);

  if (frames.length === 0) {
    return (
      <div className="preview-placeholder">
        {characterName.charAt(0).toUpperCase()}
      </div>
    );
  }

  return (
    <img
      src={frames[currentFrame]}
      alt={characterName}
      className="preview-image"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    />
  );
}

function App() {
  const [characters, setCharacters] = useState<CharacterInfo[]>([]);
  const [activeWindows, setActiveWindows] = useState<CharacterWindowInfo[]>([]);
  const [loading, setLoading] = useState(true);

  const loadCharacters = useCallback(async () => {
    setLoading(true);
    try {
      const chars = await GetCharacters();
      setCharacters(chars || []);
    } catch (err) {
      console.error('Failed to load characters:', err);
    }
    setLoading(false);
  }, []);

  const refreshActiveWindows = useCallback(async () => {
    try {
      const windows = await GetActiveWindows();
      setActiveWindows(windows || []);
    } catch (err) {
      console.error('Failed to get active windows:', err);
    }
  }, []);

  useEffect(() => {
    loadCharacters();
    const interval = setInterval(refreshActiveWindows, 1000);
    return () => clearInterval(interval);
  }, [loadCharacters, refreshActiveWindows]);

  const handleSpawn = async (characterName: string) => {
    try {
      await SpawnCharacter(characterName);
      refreshActiveWindows();
    } catch (err) {
      console.error('Failed to spawn character:', err);
    }
  };

  const handleDestroy = async (windowId: string) => {
    try {
      await DestroyCharacter(windowId);
      setTimeout(refreshActiveWindows, 100);
    } catch (err) {
      console.error('Failed to destroy character:', err);
    }
  };

  const handleScaleChange = async (windowId: string, scale: number) => {
    try {
      await SetCharacterScale(windowId, scale);
    } catch (err) {
      console.error('Failed to set scale:', err);
    }
  };

  const handleOpenFrames = async () => {
    try {
      await OpenFramesDir();
    } catch (err) {
      console.error('Failed to open frames directory:', err);
    }
  };

  const handleOpenConfig = async () => {
    try {
      await OpenConfig();
    } catch (err) {
      console.error('Failed to open config:', err);
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-left">
          <h1>Boccho Desktop</h1>
          <p className="subtitle">Character Manager</p>
        </div>
        <div className="header-right">
          <button className="btn btn-toolbar" onClick={handleOpenFrames}>
            Open Frames Dir
          </button>
          <button className="btn btn-toolbar" onClick={handleOpenConfig}>
            Open Config
          </button>
        </div>
      </header>

      <main className="app-content">
        <section className="section">
          <div className="section-header">
            <h2 className="section-title">Available Characters</h2>
            <button className="btn btn-refresh" onClick={loadCharacters}>
              Refresh
            </button>
          </div>
          {loading ? (
            <div className="loading">Loading characters...</div>
          ) : characters.length === 0 ? (
            <div className="empty-state">
              <p>No characters found</p>
              <p className="hint">
                Add character folders to the Frames directory
              </p>
            </div>
          ) : (
            <div className="character-grid">
              {characters.map((char) => (
                <div key={char.name} className="character-card">
                  <div className="character-preview">
                    <AnimatedPreview characterName={char.name} />
                  </div>
                  <div className="character-info">
                    <span className="character-name">{char.name}</span>
                    <span className="frame-count">{char.frameCount} frames</span>
                  </div>
                  <button
                    className="btn btn-spawn"
                    onClick={() => handleSpawn(char.name)}
                  >
                    Spawn
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        <section className="section">
          <h2 className="section-title">
            Active Windows
            <span className="badge">{activeWindows.length}</span>
          </h2>
          {activeWindows.length === 0 ? (
            <div className="empty-state">
              <p>No active windows</p>
              <p className="hint">Click "Spawn" to create a character window</p>
            </div>
          ) : (
            <div className="windows-list">
              {activeWindows.map((win) => (
                <div key={win.id} className="window-item">
                  <div className="window-header">
                    <div className="window-info">
                      <span className="window-name">{win.characterName}</span>
                      <span className="window-id">#{win.id}</span>
                    </div>
                    <button
                      className="btn btn-destroy"
                      onClick={() => handleDestroy(win.id)}
                    >
                      Close
                    </button>
                  </div>
                  <div className="window-controls">
                    <label className="scale-label">
                      Scale: {(win.scale * 100).toFixed(0)}%
                    </label>
                    <input
                      type="range"
                      className="scale-slider"
                      min="10"
                      max="200"
                      value={win.scale * 100}
                      onChange={(e) =>
                        handleScaleChange(win.id, parseInt(e.target.value) / 100)
                      }
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>

      <footer className="app-footer">
        <p>Arrow Keys: Resize | Escape: Close Window | Drag: Move</p>
      </footer>
    </div>
  );
}

export default App;
