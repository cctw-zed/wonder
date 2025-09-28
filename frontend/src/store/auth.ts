import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User, authAPI, LoginRequest, RegisterRequest, setAuthToken, removeAuthToken } from '@/lib/api';

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;
}

interface AuthActions {
  login: (credentials: LoginRequest) => Promise<boolean>;
  register: (userData: RegisterRequest) => Promise<boolean>;
  logout: () => void;
  clearError: () => void;
  getProfile: () => Promise<void>;
}

export const useAuthStore = create<AuthState & AuthActions>()(
  persist(
    (set, get) => ({
      // 初始状态
      user: null,
      token: null,
      isLoading: false,
      error: null,
      isAuthenticated: false,

      // 登录
      login: async (credentials: LoginRequest): Promise<boolean> => {
        set({ isLoading: true, error: null });

        try {
          console.log('🔐 Attempting login with:', { email: credentials.email });
          const response = await authAPI.login(credentials);
          console.log('✅ Login response:', response);

          const { access_token, user } = response.data;

          // 保存token到localStorage
          setAuthToken(access_token);

          // 更新状态
          set({
            user,
            token: access_token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });

          console.log('✅ Login successful, user:', user);
          return true;
        } catch (error: unknown) {
          console.error('❌ Login error:', error);
          const errorMessage = (error as { response?: { data?: { message?: string } } })?.response?.data?.message || '登录失败，请检查邮箱和密码';
          set({
            error: errorMessage,
            isLoading: false,
            isAuthenticated: false,
          });
          return false;
        }
      },

      // 注册
      register: async (userData: RegisterRequest): Promise<boolean> => {
        set({ isLoading: true, error: null });

        try {
          await authAPI.register(userData);

          // 注册成功后自动登录
          const loginSuccess = await get().login({
            email: userData.email,
            password: userData.password,
          });

          return loginSuccess;
        } catch (error: unknown) {
          const errorMessage = (error as { response?: { data?: { message?: string } } })?.response?.data?.message || '注册失败，请检查输入信息';
          set({
            error: errorMessage,
            isLoading: false,
          });
          return false;
        }
      },

      // 登出
      logout: () => {
        removeAuthToken();
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          error: null,
        });
      },

      // 清除错误
      clearError: () => {
        set({ error: null });
      },

      // 获取用户资料
      getProfile: async () => {
        if (!get().isAuthenticated) return;

        set({ isLoading: true });

        try {
          const profile = await authAPI.getProfile();
          set({
            user: profile,
            isLoading: false,
          });
        } catch (error) {
          console.error('Failed to load profile', error);
          set({ isLoading: false });
          get().logout();
        }
      },
    }),
    {
      name: 'auth-storage',
      // 只持久化必要的数据
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
