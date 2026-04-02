import type {
  GenerateConfigurationRequest,
  GenerateConfigurationResponse,
  SavedConfiguration,
} from './types';

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? '';

export async function getCurrentConfiguration(): Promise<SavedConfiguration | null> {
  const response = await fetch(`${apiBaseUrl}/api/configuration/current`);
  if (response.status === 404) {
    return null;
  }

  return handleJsonResponse<SavedConfiguration>(response);
}

export async function generateConfiguration(
  payload: GenerateConfigurationRequest,
): Promise<GenerateConfigurationResponse> {
  const response = await fetch(`${apiBaseUrl}/api/configuration/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  return handleJsonResponse<GenerateConfigurationResponse>(response);
}

export async function saveCurrentConfiguration(
  payload: Omit<SavedConfiguration, 'updatedAt'>,
): Promise<SavedConfiguration> {
  const response = await fetch(`${apiBaseUrl}/api/configuration/current`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  return handleJsonResponse<SavedConfiguration>(response);
}

async function handleJsonResponse<T>(response: Response): Promise<T> {
  const payload: unknown = await response.json();
  if (!response.ok) {
    const message =
      typeof payload === 'object' &&
      payload !== null &&
      'message' in payload &&
      typeof payload.message === 'string'
        ? payload.message
        : 'Unexpected server error.';
    throw new Error(message);
  }

  return payload as T;
}
