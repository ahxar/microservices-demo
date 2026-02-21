import axios from 'axios';
import { redirectToLogin } from '@/lib/auth-redirect';

const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    if (typeof window === 'undefined') {
      return config;
    }

    const token = window.localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const requestURL = originalRequest?.url as string | undefined;
    const isAuthRequest = Boolean(requestURL?.includes('/api/v1/auth/'));

    // If error is 401 and we haven't retried yet, try to refresh token
    if (error.response?.status === 401 && !isAuthRequest && !originalRequest?._retry) {
      originalRequest._retry = true;

      try {
        if (typeof window === 'undefined') {
          return Promise.reject(error);
        }

        const refreshToken = window.localStorage.getItem('refresh_token');
        if (!refreshToken) {
          window.localStorage.removeItem('access_token');
          window.localStorage.removeItem('refresh_token');
          redirectToLogin('auth_required');
          return Promise.reject(error);
        }

        const response = await axios.post(`${baseURL}/api/v1/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token: newRefreshToken } = response.data;
        window.localStorage.setItem('access_token', access_token);
        window.localStorage.setItem('refresh_token', newRefreshToken);

        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        if (typeof window !== 'undefined') {
          window.localStorage.removeItem('access_token');
          window.localStorage.removeItem('refresh_token');
          redirectToLogin('session_expired');
        }
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
