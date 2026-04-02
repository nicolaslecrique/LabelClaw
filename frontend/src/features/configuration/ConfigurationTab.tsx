import { useEffect, useState } from 'react';
import { GeneratedRuntime } from '../../lib/generated-runtime';
import {
  createInitialValueFromSchema,
  parseJsonInput,
  validateWithSchema,
} from '../../lib/schema';
import type { JsonValue, SavedConfiguration } from '../../lib/types';

interface ConfigurationTabProps {
  initialConfiguration: SavedConfiguration | null;
  onGenerate: (payload: {
    sampleSchema: JsonValue;
    labelSchema: JsonValue;
    uiPrompt: string;
  }) => Promise<{ componentSource: string; sampleData: JsonValue }>;
  onSave: (payload: Omit<SavedConfiguration, 'updatedAt'>) => Promise<void>;
}

export function ConfigurationTab({
  initialConfiguration,
  onGenerate,
  onSave,
}: ConfigurationTabProps) {
  const [sampleSchemaText, setSampleSchemaText] = useState(() =>
    stringifyJson(initialConfiguration?.sampleSchema ?? defaultSampleSchema),
  );
  const [labelSchemaText, setLabelSchemaText] = useState(() =>
    stringifyJson(initialConfiguration?.labelSchema ?? defaultLabelSchema),
  );
  const [uiPrompt, setUiPrompt] = useState(
    initialConfiguration?.uiPrompt ??
      'Render the sample article in a readable card with a textarea for the label.',
  );
  const [generatedCode, setGeneratedCode] = useState(
    initialConfiguration?.componentSource ?? '',
  );
  const [sampleData, setSampleData] = useState<JsonValue>(
    initialConfiguration?.sampleData ?? {},
  );
  const [labelValue, setLabelValue] = useState<JsonValue>(() =>
    createInitialValueFromSchema(initialConfiguration?.labelSchema ?? defaultLabelSchema),
  );
  const [hasLabelOutput, setHasLabelOutput] = useState(false);
  const [previewReady, setPreviewReady] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [loadMessage, setLoadMessage] = useState<string | null>(
    initialConfiguration ? 'Loaded saved configuration.' : null,
  );
  const [formError, setFormError] = useState<string | null>(null);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!initialConfiguration) {
      return;
    }

    setSampleSchemaText(stringifyJson(initialConfiguration.sampleSchema));
    setLabelSchemaText(stringifyJson(initialConfiguration.labelSchema));
    setUiPrompt(initialConfiguration.uiPrompt);
    setGeneratedCode(initialConfiguration.componentSource);
    setSampleData(initialConfiguration.sampleData);
    setLabelValue(createInitialValueFromSchema(initialConfiguration.labelSchema));
    setHasLabelOutput(false);
    setPreviewReady(false);
    setLoadMessage('Loaded saved configuration.');
  }, [initialConfiguration]);

  const sampleSchema =
    safeParseJson(sampleSchemaText, 'Sample JSON Schema') ?? null;
  const labelSchema = safeParseJson(labelSchemaText, 'Label JSON Schema') ?? null;
  const sampleSchemaError =
    sampleSchema === null ? 'Sample JSON Schema must be valid JSON.' : null;
  const labelSchemaError =
    labelSchema === null ? 'Label JSON Schema must be valid JSON.' : null;
  const labelValidationError =
    hasLabelOutput && labelSchema
      ? validateWithSchema(labelSchema, labelValue)
      : null;

  async function handleGenerate() {
    if (!sampleSchema || !labelSchema) {
      setFormError('Both schemas must be valid JSON before generation.');
      return;
    }

    if (!uiPrompt.trim()) {
      setFormError('UI prompt must not be empty.');
      return;
    }

    setIsGenerating(true);
    setFormError(null);
    setSaveMessage(null);
    setLoadMessage(null);

    try {
      const response = await onGenerate({
        sampleSchema,
        labelSchema,
        uiPrompt,
      });
      setGeneratedCode(response.componentSource);
      setSampleData(response.sampleData);
      setLabelValue(createInitialValueFromSchema(labelSchema));
      setHasLabelOutput(false);
      setPreviewReady(false);
    } catch (error) {
      setFormError(
        error instanceof Error ? error.message : 'Failed to generate panel.',
      );
    } finally {
      setIsGenerating(false);
    }
  }

  async function handleSave() {
    if (!sampleSchema || !labelSchema) {
      setFormError('Both schemas must be valid JSON before saving.');
      return;
    }

    setIsSaving(true);
    setFormError(null);
    setSaveMessage(null);

    try {
      await onSave({
        sampleSchema,
        labelSchema,
        uiPrompt,
        sampleData,
        componentSource: generatedCode,
      });
      setSaveMessage('Labelling panel saved.');
    } catch (error) {
      setFormError(error instanceof Error ? error.message : 'Save failed.');
    } finally {
      setIsSaving(false);
    }
  }

  return (
    <section className="workspace-grid">
      <div className="editor-panel">
        <p className="eyebrow">Configuration</p>
        <h2>Design the labelling panel</h2>
        <p className="panel-copy">
          Define the sample and label contracts, describe the UI, then generate
          a React panel from the backend.
        </p>

        <label className="field">
          <span>Sample JSON Schema</span>
          <textarea
            aria-label="Sample JSON Schema"
            value={sampleSchemaText}
            onChange={(event) => setSampleSchemaText(event.target.value)}
            rows={12}
          />
          {sampleSchemaError ? <small>{sampleSchemaError}</small> : null}
        </label>

        <label className="field">
          <span>Label JSON Schema</span>
          <textarea
            aria-label="Label JSON Schema"
            value={labelSchemaText}
            onChange={(event) => setLabelSchemaText(event.target.value)}
            rows={12}
          />
          {labelSchemaError ? <small>{labelSchemaError}</small> : null}
        </label>

        <label className="field">
          <span>UI prompt</span>
          <textarea
            aria-label="UI prompt"
            value={uiPrompt}
            onChange={(event) => setUiPrompt(event.target.value)}
            rows={7}
          />
        </label>

        <button
          className="primary-button"
          type="button"
          onClick={() => {
            void handleGenerate();
          }}
          disabled={isGenerating}
        >
          {isGenerating ? 'Generating…' : 'Generate labelling panel'}
        </button>

        {formError ? (
          <p className="message message-error" role="alert">
            {formError}
          </p>
        ) : null}
        {loadMessage ? <p className="message">{loadMessage}</p> : null}
      </div>

      <div className="preview-panel">
        <div className="preview-toolbar">
          <div>
            <p className="eyebrow">Preview</p>
            <h2>Generated labelling panel</h2>
          </div>
          <button
            className="secondary-button"
            type="button"
            onClick={() => {
              void handleSave();
            }}
            disabled={!generatedCode || !previewReady || isSaving}
          >
            {isSaving ? 'Saving…' : 'Save labelling panel'}
          </button>
        </div>

        <div className="preview-card">
          <GeneratedRuntime
            code={generatedCode}
            sample={sampleData}
            value={labelValue}
            onChange={(nextValue) => {
              setLabelValue(nextValue);
              setHasLabelOutput(true);
            }}
            onReadyChange={setPreviewReady}
          />
        </div>

        <section className="inspection-grid">
          <article className="inspection-panel">
            <h3>Sample data</h3>
            <pre>{stringifyJson(sampleData)}</pre>
          </article>
          <article className="inspection-panel">
            <h3>Current label output</h3>
            {hasLabelOutput ? (
              <pre>{stringifyJson(labelValue)}</pre>
            ) : (
              <p>The generated component has not emitted a label yet.</p>
            )}
            {labelValidationError ? (
              <p className="message message-error" role="alert">
                Output mismatch: {labelValidationError}
              </p>
            ) : hasLabelOutput ? (
              <p className="message message-success">Output matches the label schema.</p>
            ) : null}
          </article>
        </section>

        {saveMessage ? <p className="message message-success">{saveMessage}</p> : null}
      </div>
    </section>
  );
}

function safeParseJson(text: string, fieldName: string): JsonValue | null {
  try {
    return parseJsonInput(text, fieldName);
  } catch {
    return null;
  }
}

function stringifyJson(value: JsonValue): string {
  return JSON.stringify(value, null, 2);
}

const defaultSampleSchema: JsonValue = {
  type: 'object',
  properties: {
    title: { type: 'string' },
    body: { type: 'string' },
  },
  required: ['title', 'body'],
  additionalProperties: false,
};

const defaultLabelSchema: JsonValue = {
  type: 'object',
  properties: {
    sentiment: {
      type: 'string',
      enum: ['positive', 'neutral', 'negative'],
    },
    notes: { type: 'string' },
  },
  required: ['sentiment', 'notes'],
  additionalProperties: false,
};
