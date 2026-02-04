/*
App.tsx - Main character manager UI

Components:
- App: Main application component with character grid and active windows list
- AnimatedPreview: Component that cycles through frames for animation preview
- AddDropdown: Dropdown menu for adding packs (from Link or .bfk)
- AddPackModal: Modal for confirming pack installation with preview
*/

import { useState, useEffect, useCallback, useRef } from 'react';
import './App.css';
import { CharacterInfo, CharacterWindowInfo, PackInfo } from './types';
import {
  GetCharacters,
  SpawnCharacter,
  DestroyCharacter,
  GetActiveWindows,
  SetCharacterScale,
  GetPreviewFrames,
  OpenFramesDir,
  OpenConfig,
  BrowseBfkFile,
  GetBfkPackInfo,
  InstallBfkPack,
} from '../wailsjs/go/main/App';

import linkIcon from './assets/images/link.svg';
import plusIcon from './assets/images/Plus_button.svg';

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

interface AddDropdownProps {
  onAddFromFile: () => void;
  onAddFromLink: () => void;
}

function AddDropdown({ onAddFromFile, onAddFromLink }: AddDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="add-dropdown" ref={dropdownRef}>
      <button
        className="btn btn-add"
        onClick={() => setIsOpen(!isOpen)}
      >
        import
      </button>
      {isOpen && (
        <div className="add-dropdown-menu">
          <button
            className="add-dropdown-item"
            onClick={() => {
              onAddFromLink();
              setIsOpen(false);
            }}
          >
            <div className="add-dropdown-icon">
              <img src={linkIcon} alt="" />
            </div>
            <span>From Link</span>
          </button>
          <button
            className="add-dropdown-item"
            onClick={() => {
              onAddFromFile();
              setIsOpen(false);
            }}
          >
            <div className="add-dropdown-icon">
              <img src={plusIcon} alt="" />
            </div>
            <span>From .bfk</span>
          </button>
        </div>
      )}
    </div>
  );
}

interface AddPackModalProps {
  packInfo: PackInfo;
  onInstall: () => void;
  onCancel: () => void;
  installing: boolean;
}

function AddPackModal({ packInfo, onInstall, onCancel, installing }: AddPackModalProps) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3 className="modal-title">Install Pack</h3>
          <p className="modal-subtitle">{packInfo.packName}</p>
        </div>

        {packInfo.previewImage && (
          <div className="modal-preview">
            <img src={packInfo.previewImage} alt="Preview" />
          </div>
        )}

        {packInfo.error ? (
          <p className="modal-error">{packInfo.error}</p>
        ) : (
          <p className="modal-info">
            Characters: <span>{packInfo.characters.join(', ')}</span>
          </p>
        )}

        <div className="modal-actions">
          <button className="btn btn-cancel" onClick={onCancel} disabled={installing}>
            Cancel
          </button>
          {!packInfo.error && (
            <button className="btn btn-install" onClick={onInstall} disabled={installing}>
              {installing ? 'Installing...' : 'Install'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

function App() {
  const [characters, setCharacters] = useState<CharacterInfo[]>([]);
  const [activeWindows, setActiveWindows] = useState<CharacterWindowInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [packInfo, setPackInfo] = useState<PackInfo | null>(null);
  const [installing, setInstalling] = useState(false);

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

  const handleAddFromFile = async () => {
    try {
      const filePath = await BrowseBfkFile();
      if (!filePath) return;

      const info = await GetBfkPackInfo(filePath);
      setPackInfo(info);
    } catch (err) {
      console.error('Failed to browse pack:', err);
    }
  };

  const handleAddFromLink = () => {
    // TODO: Implement add from link
    console.log('Add from link - not implemented yet');
  };

  const handleInstallPack = async () => {
    if (!packInfo) return;

    setInstalling(true);
    try {
      await InstallBfkPack(packInfo.filePath);
      setPackInfo(null);
      loadCharacters();
    } catch (err) {
      console.error('Failed to install pack:', err);
    }
    setInstalling(false);
  };

  const handleCancelPack = () => {
    setPackInfo(null);
  };

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-left">
          <h1>Boccho Desktop</h1>
          <p className="subtitle">Character Manager</p>
        </div>
        <div className="header-right">
          <AddDropdown
            onAddFromFile={handleAddFromFile}
            onAddFromLink={handleAddFromLink}
          />
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

      {packInfo && (
        <AddPackModal
          packInfo={packInfo}
          onInstall={handleInstallPack}
          onCancel={handleCancelPack}
          installing={installing}
        />
      )}
    </div>
  );
}

export default App;
