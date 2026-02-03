/*
App.tsx - Main character manager UI

Components:
- App: Main application component with character grid and active windows list
- Character cards with preview images
- Active windows with scale slider
*/

import { useState, useEffect, useCallback } from 'react';
import './App.css';
import { CharacterInfo, CharacterWindowInfo } from './types';
import {
  GetCharacters,
  SpawnCharacter,
  DestroyCharacter,
  GetActiveWindows,
  SetCharacterScale,
  GetPreviewImageBase64,
} from '../wailsjs/go/main/App';

function App() {
  const [characters, setCharacters] = useState<CharacterInfo[]>([]);
  const [previews, setPreviews] = useState<Record<string, string>>({});
  const [activeWindows, setActiveWindows] = useState<CharacterWindowInfo[]>([]);
  const [loading, setLoading] = useState(true);

  const loadCharacters = useCallback(async () => {
    try {
      const chars = await GetCharacters();
      setCharacters(chars || []);

      const previewPromises = (chars || []).map(async (char) => {
        const preview = await GetPreviewImageBase64(char.name);
        return { name: char.name, preview };
      });

      const previewResults = await Promise.all(previewPromises);
      const previewMap: Record<string, string> = {};
      previewResults.forEach(({ name, preview }) => {
        if (preview) previewMap[name] = preview;
      });
      setPreviews(previewMap);
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

  return (
    <div className="app">
      <header className="app-header">
        <h1>Boccho Desktop</h1>
        <p className="subtitle">Character Manager</p>
      </header>

      <main className="app-content">
        <section className="section">
          <h2 className="section-title">Available Characters</h2>
          {loading ? (
            <div className="loading">Loading characters...</div>
          ) : characters.length === 0 ? (
            <div className="empty-state">
              <p>No characters found</p>
              <p className="hint">Add character folders to the Frames directory</p>
            </div>
          ) : (
            <div className="character-grid">
              {characters.map((char) => (
                <div key={char.name} className="character-card">
                  <div className="character-preview">
                    {previews[char.name] ? (
                      <img
                        src={previews[char.name]}
                        alt={char.name}
                        className="preview-image"
                      />
                    ) : (
                      <div className="preview-placeholder">
                        {char.name.charAt(0).toUpperCase()}
                      </div>
                    )}
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
