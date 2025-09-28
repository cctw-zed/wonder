import axios from 'axios';

// API configuration
const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const API_VERSION = '/api/v1';

console.log('🔧 API Configuration:', { BASE_URL, API_VERSION, fullURL: `${BASE_URL}${API_VERSION}` });

// Axios instance shared by all API helpers
export const api = axios.create({
  baseURL: `${BASE_URL}${API_VERSION}`,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

// Attach auth token before every request
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Handle auth failures globally
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Remove any persisted auth data when the session expires
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth-storage');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Domain models
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponseData {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface AuthResponse {
  data: AuthResponseData;
  trace_id: string;
}

interface RawMeResponse {
  user_id?: string;
  data?: Partial<User> & { id?: string };
  trace_id: string;
}

export interface UpdateProfilePayload {
  email?: string;
  name?: string;
}

export interface ChangePasswordPayload {
  old_password: string;
  new_password: string;
}

const resolveUserId = (response: RawMeResponse): string | undefined => {
  if (response.user_id) {
    return response.user_id;
  }

  const candidateId = response.data?.id;
  if (candidateId) {
    return candidateId;
  }

  return undefined;
};

// API helpers grouped by concern
export const authAPI = {
  register: async (userData: RegisterRequest): Promise<{ user: User; trace_id: string }> => {
    const response = await api.post('/users/register', userData);
    return response.data;
  },

  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', credentials);
    return response.data;
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get<RawMeResponse>('/auth/me');
    const directUser = response.data.data;

    if (directUser && directUser.email && directUser.name && directUser.id) {
      return directUser as User;
    }

    const userId = resolveUserId(response.data);
    if (!userId) {
      throw new Error('Unable to resolve authenticated user id from /auth/me response');
    }

    const profileResponse = await api.get<{ user: User }>(`/users/${userId}`);
    const user = profileResponse.data.user;

    if (!user) {
      throw new Error('User profile response did not include user data');
    }

    return user;
  },
};

export const userAPI = {
  updateProfile: async (userId: string, payload: UpdateProfilePayload): Promise<User> => {
    const response = await api.put<{ user: User }>(`/users/${userId}`, payload);
    return response.data.user;
  },

  changePassword: async (userId: string, payload: ChangePasswordPayload): Promise<void> => {
    await api.put(`/users/${userId}/password`, payload);
  },
};

// Auth token helpers
export const setAuthToken = (token: string) => {
  localStorage.setItem('auth_token', token);
};

export const getAuthToken = (): string | null => {
  return localStorage.getItem('auth_token');
};

export const removeAuthToken = () => {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('auth-storage');
};

export const isAuthenticated = (): boolean => {
  return !!getAuthToken();
};
