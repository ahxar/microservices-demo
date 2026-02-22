type ErrorWithResponse = {
  message?: string;
  response?: {
    data?: {
      message?: string;
      error?: string;
    };
  };
};

export function getErrorMessage(error: unknown, fallback: string): string {
  if (!error) {
    return fallback;
  }

  if (typeof error === 'string' && error.trim().length > 0) {
    return error;
  }

  const err = error as ErrorWithResponse;
  return err.response?.data?.message || err.response?.data?.error || err.message || fallback;
}
