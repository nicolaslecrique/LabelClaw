import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { App } from './App';

const fetchMock = vi.fn<typeof fetch>();

describe('App', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', fetchMock);
    fetchMock.mockReset();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    cleanup();
  });

  it('shows inline validation when schemas are invalid', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ message: 'No saved configuration found.' }), {
        status: 404,
        headers: { 'Content-Type': 'application/json' },
      }),
    );

    const user = userEvent.setup();
    render(<App />);

    fireEvent.change(screen.getByLabelText('Sample JSON Schema'), {
      target: { value: '{' },
    });
    await user.click(screen.getByRole('button', { name: 'Generate labelling panel' }));

    expect(
      await screen.findByText('Both schemas must be valid JSON before generation.'),
    ).toBeInTheDocument();
  });

  it('generates a preview and enables saving', async () => {
    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ message: 'No saved configuration found.' }), {
          status: 404,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            componentSource:
              "export default function LabelingPanel({ sample, value, onChange }) { return <div><p>{sample.title}</p><textarea aria-label='Generated notes' value={typeof value === 'string' ? value : ''} onChange={(event) => onChange(event.target.value)} /></div>; }",
            sampleData: { title: 'Hello world' },
          }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            sampleSchema: { type: 'object' },
            labelSchema: { type: 'object' },
            uiPrompt: 'Prompt',
            sampleData: { title: 'Hello world' },
            componentSource:
              "export default function LabelingPanel({ sample, value, onChange }) { return <div><p>{sample.title}</p><textarea aria-label='Generated notes' value={typeof value === 'string' ? value : ''} onChange={(event) => onChange(event.target.value)} /></div>; }",
            updatedAt: '2026-04-02T10:30:00Z',
          }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      );

    const user = userEvent.setup();
    render(<App />);

    fireEvent.change(screen.getByLabelText('Label JSON Schema'), {
      target: { value: JSON.stringify({ type: 'string' }) },
    });
    await user.click(screen.getByRole('button', { name: 'Generate labelling panel' }));

    expect(await screen.findByText('Hello world')).toBeInTheDocument();

    const saveButton = screen.getByRole('button', { name: 'Save labelling panel' });
    await waitFor(() => expect(saveButton).toBeEnabled());
    await user.click(saveButton);

    expect(await screen.findByText('Labelling panel saved.')).toBeInTheDocument();
  });

  it('shows preview errors without crashing the page', async () => {
    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ message: 'No saved configuration found.' }), {
          status: 404,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            componentSource:
              "export default function LabelingPanel() { throw new Error('boom'); }",
            sampleData: { title: 'Hello world' },
          }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      );

    const user = userEvent.setup();
    render(<App />);

    fireEvent.change(screen.getByLabelText('Label JSON Schema'), {
      target: { value: JSON.stringify({ type: 'string' }) },
    });
    await user.click(screen.getByRole('button', { name: 'Generate labelling panel' }));

    expect(await screen.findByText('Preview failed')).toBeInTheDocument();
    expect(await screen.findByText('boom')).toBeInTheDocument();
  });

  it('surfaces label schema mismatches after the generated UI emits data', async () => {
    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ message: 'No saved configuration found.' }), {
          status: 404,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            componentSource:
              "export default function LabelingPanel({ onChange }) { return <button type='button' onClick={() => onChange(123)}>Emit</button>; }",
            sampleData: { title: 'Hello world' },
          }),
          {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          },
        ),
      );

    const user = userEvent.setup();
    render(<App />);

    fireEvent.change(screen.getByLabelText('Label JSON Schema'), {
      target: { value: JSON.stringify({ type: 'string' }) },
    });
    await user.click(screen.getByRole('button', { name: 'Generate labelling panel' }));
    await user.click(await screen.findByRole('button', { name: 'Emit' }));

    expect(await screen.findByText(/Output mismatch:/)).toBeInTheDocument();
  });

  it('keeps the labelling tab as a placeholder', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ message: 'No saved configuration found.' }), {
        status: 404,
        headers: { 'Content-Type': 'application/json' },
      }),
    );

    const user = userEvent.setup();
    render(<App />);

    await user.click(screen.getByRole('button', { name: 'Labelling' }));

    expect(await screen.findByText('To be defined later')).toBeInTheDocument();
  });
});
