import { createServer } from 'node:http';

const port = 8080;

let savedConfiguration = null;

const generatedPanel = {
  componentSource:
    "export default function LabelingPanel({ sample, value, onChange }) { return <section><h3>{sample.title}</h3><p>{sample.body}</p><select aria-label='Sentiment' value={value?.sentiment ?? ''} onChange={(event) => onChange({ ...(value && typeof value === 'object' ? value : {}), sentiment: event.target.value, notes: value?.notes ?? '' })}><option value=''>Choose</option><option value='positive'>Positive</option><option value='neutral'>Neutral</option><option value='negative'>Negative</option></select><textarea aria-label='Notes' value={value?.notes ?? ''} onChange={(event) => onChange({ ...(value && typeof value === 'object' ? value : {}), sentiment: value?.sentiment ?? '', notes: event.target.value })} /></section>; }",
  sampleData: {
    title: 'Example article',
    body: 'A short sample item for previewing the generated panel.',
  },
};

const server = createServer(async (request, response) => {
  response.setHeader('Access-Control-Allow-Origin', '*');
  response.setHeader('Access-Control-Allow-Headers', 'Content-Type');
  response.setHeader('Access-Control-Allow-Methods', 'GET,POST,PUT,OPTIONS');

  if (request.method === 'OPTIONS') {
    response.writeHead(204);
    response.end();
    return;
  }

  if (request.url === '/api/health' && request.method === 'GET') {
    writeJson(response, 200, { status: 'ok' });
    return;
  }

  if (request.url === '/api/configuration/current' && request.method === 'GET') {
    if (!savedConfiguration) {
      writeJson(response, 404, { message: 'No saved configuration found.' });
      return;
    }

    writeJson(response, 200, savedConfiguration);
    return;
  }

  if (request.url === '/api/configuration/generate' && request.method === 'POST') {
    const body = await readJson(request);

    if (!body.uiPrompt || !body.sampleSchema || !body.labelSchema) {
      writeJson(response, 400, { message: 'Invalid generation payload.' });
      return;
    }

    writeJson(response, 200, generatedPanel);
    return;
  }

  if (request.url === '/api/configuration/current' && request.method === 'PUT') {
    const body = await readJson(request);
    savedConfiguration = {
      ...body,
      updatedAt: new Date().toISOString(),
    };
    writeJson(response, 200, savedConfiguration);
    return;
  }

  writeJson(response, 404, { message: 'Not found.' });
});

server.listen(port, '127.0.0.1', () => {
  process.stdout.write(`e2e stub server listening on http://127.0.0.1:${port}\n`);
});

function writeJson(response, statusCode, payload) {
  response.writeHead(statusCode, {
    'Content-Type': 'application/json',
  });
  response.end(JSON.stringify(payload));
}

async function readJson(request) {
  const chunks = [];
  for await (const chunk of request) {
    chunks.push(chunk);
  }

  return JSON.parse(Buffer.concat(chunks).toString('utf8'));
}

