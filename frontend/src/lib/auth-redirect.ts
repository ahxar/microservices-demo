const LOGIN_PATH = '/login';

export type LoginRedirectReason = 'auth_required' | 'action_requires_auth' | 'session_expired';

export function hasAccessToken(): boolean {
  if (typeof window === 'undefined') {
    return false;
  }

  return Boolean(window.localStorage.getItem('access_token'));
}

export function redirectToLogin(reason: LoginRedirectReason): void {
  if (typeof window === 'undefined') {
    return;
  }

  const params = new URLSearchParams();
  params.set('reason', reason);

  const next = `${window.location.pathname}${window.location.search}`;
  if (next && next !== LOGIN_PATH) {
    params.set('next', next);
  }

  window.location.href = `${LOGIN_PATH}?${params.toString()}`;
}

export function getLoginReasonMessage(reason: string | null): string | null {
  switch (reason) {
    case 'action_requires_auth':
      return 'Please log in to add items to cart or wishlist.';
    case 'session_expired':
      return 'Your session expired. Please log in again.';
    case 'auth_required':
      return 'Please log in to continue.';
    default:
      return null;
  }
}

export function getSafeNextPath(nextPath: string | null): string {
  if (!nextPath || !nextPath.startsWith('/')) {
    return '/products';
  }

  return nextPath;
}
