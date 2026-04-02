import Ajv2020, {
  type AnySchema,
  type ErrorObject,
  type ValidateFunction,
} from 'ajv/dist/2020';
import type { JsonValue } from './types';

const ajv = new Ajv2020({
  allErrors: true,
  strict: false,
});

export function parseJsonInput(text: string, fieldName: string): JsonValue {
  try {
    return JSON.parse(text) as JsonValue;
  } catch {
    throw new Error(`${fieldName} must be valid JSON.`);
  }
}

export function validateWithSchema(schema: JsonValue, value: JsonValue): string | null {
  const validator = compileSchema(schema);
  if (validator(value)) {
    return null;
  }

  return formatAjvError(validator.errors);
}

export function createInitialValueFromSchema(schema: JsonValue): JsonValue {
  if (!schema || typeof schema !== 'object' || Array.isArray(schema)) {
    return null;
  }

  if ('default' in schema) {
    return schema.default;
  }

  switch (schema.type) {
    case 'object':
      return {};
    case 'array':
      return [];
    case 'boolean':
      return false;
    case 'integer':
    case 'number':
      return 0;
    case 'string':
      return '';
    case 'null':
      return null;
    default:
      return null;
  }
}

function compileSchema(schema: JsonValue): ValidateFunction<JsonValue> {
  if (!schema || typeof schema !== 'object' || Array.isArray(schema)) {
    throw new Error('Schema must be a JSON object.');
  }

  return ajv.compile(schema as AnySchema);
}

function formatAjvError(errors: ErrorObject[] | null | undefined): string {
  if (!errors || errors.length === 0) {
    return 'Value does not match the schema.';
  }

  const [firstError] = errors;
  if (!firstError) {
    return 'Value does not match the schema.';
  }

  const location = firstError.instancePath || 'value';
  return `${location} ${firstError.message ?? 'is invalid'}`.trim();
}
