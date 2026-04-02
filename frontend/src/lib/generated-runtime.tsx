import { transform } from '@babel/standalone';
import {
  Component,
  type ReactNode,
  useEffect,
  useEffectEvent,
  useState,
} from 'react';
import * as React from 'react';
import type { JsonValue } from './types';

export interface GeneratedComponentProps {
  sample: JsonValue;
  value: JsonValue;
  onChange: (nextValue: JsonValue) => void;
}

interface GeneratedRuntimeProps extends GeneratedComponentProps {
  code: string;
  onReadyChange?: (isReady: boolean) => void;
}

export function GeneratedRuntime({
  code,
  sample,
  value,
  onChange,
  onReadyChange,
}: GeneratedRuntimeProps) {
  const [GeneratedComponent, setGeneratedComponent] =
    useState<React.ComponentType<GeneratedComponentProps> | null>(null);
  const [error, setError] = useState<string | null>(null);
  const notifyReadyChange = useEffectEvent((isReady: boolean) => {
    onReadyChange?.(isReady);
  });

  useEffect(() => {
    if (!code.trim()) {
      setGeneratedComponent(null);
      setError(null);
      return;
    }

    try {
      const nextComponent = compileGeneratedComponent(code);
      setGeneratedComponent(() => nextComponent);
      setError(null);
    } catch (runtimeError) {
      setGeneratedComponent(null);
      setError(
        runtimeError instanceof Error
          ? runtimeError.message
          : 'Failed to compile the generated component.',
      );
    }
  }, [code]);

  useEffect(() => {
    if (error) {
      notifyReadyChange(false);
    }
  }, [error, notifyReadyChange]);

  if (error) {
    return (
      <section className="preview-state preview-state-error" role="alert">
        <h3>Preview failed</h3>
        <p>{error}</p>
      </section>
    );
  }

  if (!GeneratedComponent) {
    return (
      <section className="preview-state">
        <h3>No panel generated yet</h3>
        <p>Generate a labelling panel to preview it here.</p>
      </section>
    );
  }

  return (
    <PreviewErrorBoundary
      key={code}
      onError={(runtimeError) => {
        setError(runtimeError.message);
        notifyReadyChange(false);
      }}
    >
      <MountedSignal onMount={() => notifyReadyChange(true)} />
      <GeneratedComponent sample={sample} value={value} onChange={onChange} />
    </PreviewErrorBoundary>
  );
}

class PreviewErrorBoundary extends Component<
  { children: ReactNode; onError: (error: Error) => void },
  { error: Error | null }
> {
  constructor(props: { children: ReactNode; onError: (error: Error) => void }) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  componentDidCatch(error: Error) {
    this.props.onError(error);
  }

  render() {
    if (this.state.error) {
      return (
        <section className="preview-state preview-state-error" role="alert">
          <h3>Preview crashed</h3>
          <p>{this.state.error.message}</p>
        </section>
      );
    }

    return this.props.children;
  }
}

function MountedSignal({ onMount }: { onMount: () => void }) {
  useEffect(() => {
    onMount();
  }, [onMount]);

  return null;
}

function compileGeneratedComponent(code: string): React.ComponentType<GeneratedComponentProps> {
  if (!code.trim()) {
    throw new Error('No generated component code is available yet.');
  }

  const componentName = extractComponentName(code);
  const normalizedSource = code.replace(
    /export\s+default\s+function\s+([A-Za-z_$][\w$]*)\s*\(/,
    'function $1(',
  );
  const compiled = transform(normalizedSource, {
    presets: ['react'],
    filename: 'generated-panel.jsx',
  }).code;

  if (!compiled) {
    throw new Error('Generated component compilation produced no code.');
  }

  const factoryBody = `${compiled}\nreturn ${componentName};`;
  // eslint-disable-next-line @typescript-eslint/no-implied-eval -- Trusted generated UI is executed in-process in v1 by design.
  const factory = new Function('React', factoryBody) as (
    react: typeof React,
  ) => React.ComponentType<GeneratedComponentProps>;

  const candidate = factory(React);
  if (typeof candidate !== 'function') {
    throw new Error('Generated component did not evaluate to a React component.');
  }

  return candidate;
}

function extractComponentName(code: string): string {
  if (/\bimport\b/.test(code) || /\brequire\s*\(/.test(code)) {
    throw new Error('Generated components cannot import modules.');
  }

  const match = code.match(
    /export\s+default\s+function\s+([A-Za-z_$][\w$]*)\s*\(/,
  );
  if (!match) {
    throw new Error(
      'Generated component must use `export default function ComponentName(...)`.',
    );
  }

  return match[1] ?? '';
}
