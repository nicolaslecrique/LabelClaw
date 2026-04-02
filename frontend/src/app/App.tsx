import { useEffect, useState, startTransition } from 'react';
import { ConfigurationTab } from '../features/configuration/ConfigurationTab';
import { LabellingTab } from '../features/labelling/LabellingTab';
import {
  generateConfiguration,
  getCurrentConfiguration,
  saveCurrentConfiguration,
} from '../lib/api';
import type { JsonValue, SavedConfiguration } from '../lib/types';

type TabId = 'configuration' | 'labelling';

export function App() {
  const [activeTab, setActiveTab] = useState<TabId>('configuration');
  const [currentConfiguration, setCurrentConfiguration] =
    useState<SavedConfiguration | null>(null);
  const [loadingError, setLoadingError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function loadCurrentConfiguration() {
      try {
        const configuration = await getCurrentConfiguration();
        if (cancelled) {
          return;
        }

        startTransition(() => {
          setCurrentConfiguration(configuration);
        });
      } catch (error) {
        if (cancelled) {
          return;
        }

        setLoadingError(
          error instanceof Error
            ? error.message
            : 'Failed to load the current configuration.',
        );
      }
    }

    void loadCurrentConfiguration();

    return () => {
      cancelled = true;
    };
  }, []);

  async function handleGenerate(payload: {
    sampleSchema: JsonValue;
    labelSchema: JsonValue;
    uiPrompt: string;
  }) {
    return generateConfiguration(payload);
  }

  async function handleSave(payload: Omit<SavedConfiguration, 'updatedAt'>) {
    const saved = await saveCurrentConfiguration(payload);
    setCurrentConfiguration(saved);
  }

  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">LabelClaw</p>
        <h1>Compose the labelling UI before the first annotation ever happens.</h1>
        <p className="hero-copy">
          Define input and output schemas, generate a task-specific React panel
          from the backend, preview it with sample data, and save the active
          configuration for future use.
        </p>
      </section>

      <section className="tabs-panel">
        <nav aria-label="Main tabs" className="tabs">
          <button
            className={activeTab === 'configuration' ? 'tab tab-active' : 'tab'}
            type="button"
            onClick={() => setActiveTab('configuration')}
          >
            Configuration
          </button>
          <button
            className={activeTab === 'labelling' ? 'tab tab-active' : 'tab'}
            type="button"
            onClick={() => setActiveTab('labelling')}
          >
            Labelling
          </button>
        </nav>

        {loadingError ? (
          <p className="message message-error" role="alert">
            {loadingError}
          </p>
        ) : null}

        {activeTab === 'configuration' ? (
          <ConfigurationTab
            initialConfiguration={currentConfiguration}
            onGenerate={handleGenerate}
            onSave={handleSave}
          />
        ) : (
          <LabellingTab />
        )}
      </section>
    </main>
  );
}

