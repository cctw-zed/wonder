'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/store/auth';
import MonitoringLinks from '@/components/monitoring/MonitoringLinks';
import MonitoringDropdown from '@/components/monitoring/MonitoringDropdown';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { userAPI } from '@/lib/api';

const profileSchema = z.object({
  name: z
    .string()
    .min(2, '姓名至少需要2个字符')
    .max(50, '姓名不能超过50个字符'),
  email: z
    .string()
    .min(1, '请输入邮箱地址')
    .email('请输入有效的邮箱地址'),
});

const passwordSchema = z
  .object({
    oldPassword: z
      .string()
      .min(6, '旧密码至少6位字符')
      .max(50, '旧密码不能超过50位字符'),
    newPassword: z
      .string()
      .min(8, '新密码至少8位字符')
      .max(50, '新密码不能超过50位字符'),
    confirmPassword: z
      .string()
      .min(8, '确认密码至少8位字符')
      .max(50, '确认密码不能超过50位字符'),
  })
  .refine((values) => values.newPassword === values.confirmPassword, {
    path: ['confirmPassword'],
    message: '两次输入的密码不一致',
  });

export default function DashboardPage() {
  const router = useRouter();
  const { user, logout, isAuthenticated, getProfile, isLoading } = useAuthStore();
  const [isProfileModalOpen, setIsProfileModalOpen] = useState(false);
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);
  const [profileStatus, setProfileStatus] = useState<{ type: 'success' | 'error'; message: string } | null>(null);
  const [passwordStatus, setPasswordStatus] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const profileForm = useForm<z.infer<typeof profileSchema>>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name ?? '',
      email: user?.email ?? '',
    },
  });

  const passwordForm = useForm<z.infer<typeof passwordSchema>>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  });

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }

    // 获取最新的用户信息
    getProfile();
  }, [isAuthenticated, router, getProfile]);

  useEffect(() => {
    if (user) {
      profileForm.reset({
        name: user.name,
        email: user.email,
      });
    }
  }, [user, profileForm]);

  useEffect(() => {
    if (!isProfileModalOpen) {
      setProfileStatus(null);
    }
  }, [isProfileModalOpen]);

  useEffect(() => {
    if (!isPasswordModalOpen) {
      setPasswordStatus(null);
      passwordForm.reset();
    }
  }, [isPasswordModalOpen, passwordForm]);

  const handleProfileSubmit = profileForm.handleSubmit(async (values) => {
    if (!user) return;

    setProfileStatus(null);

    try {
      const updatedUser = await userAPI.updateProfile(user.id, values);
      useAuthStore.setState({ user: updatedUser });
      setProfileStatus({ type: 'success', message: '个人资料已更新' });
    } catch (error: unknown) {
      const message =
        (error as { response?: { data?: { message?: string } } })?.response?.data?.message || '更新个人资料失败，请稍后重试';
      setProfileStatus({ type: 'error', message });
    }
  });

  const handlePasswordSubmit = passwordForm.handleSubmit(async (values) => {
    if (!user) return;

    setPasswordStatus(null);

    try {
      await userAPI.changePassword(user.id, {
        old_password: values.oldPassword,
        new_password: values.newPassword,
      });
      setPasswordStatus({ type: 'success', message: '密码更新成功' });
      passwordForm.reset();
    } catch (error: unknown) {
      const message =
        (error as { response?: { data?: { message?: string } } })?.response?.data?.message || '更新密码失败，请稍后重试';
      setPasswordStatus({ type: 'error', message });
    }
  });

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  if (!isAuthenticated) {
    return null; // 重定向到登录页面
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-dashboard">
      {/* 顶部导航 */}
      <header className="glass-card border-0 shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center space-x-4">
              <div className="icon-container w-10 h-10">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900">Wonder</h1>
            </div>
            <div className="flex items-center space-x-4">
              <div className="hidden sm:flex items-center space-x-3 bg-white/60 backdrop-blur-sm rounded-xl px-4 py-2">
                <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-sm font-semibold">
                  {user?.name?.charAt(0).toUpperCase()}
                </div>
                <div className="text-sm">
                  <p className="font-medium text-gray-900">欢迎, {user?.name}</p>
                  <p className="text-gray-500">{user?.email}</p>
                </div>
              </div>
              <MonitoringDropdown />
              <button
                onClick={handleLogout}
                className="btn-secondary text-red-600 border-red-200 hover:bg-red-50"
              >
                <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                退出登录
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto py-8 sm:px-6 lg:px-8">
        <div className="px-4 sm:px-0">
          <div className="space-y-8">
            {/* 欢迎卡片 */}
            <div className="glass-card rounded-2xl p-8 card-hover">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">
                    欢迎来到您的仪表板
                  </h2>
                  <p className="text-gray-600 mb-6">
                    这里是您的个人工作区，您可以查看和管理您的信息
                  </p>

                  <div className="bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200/50 rounded-xl p-4">
                    <div className="flex items-center">
                      <div className="flex-shrink-0">
                        <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center">
                          <svg className="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                          </svg>
                        </div>
                      </div>
                      <div className="ml-4">
                        <h3 className="text-sm font-semibold text-green-800">登录成功！</h3>
                        <p className="text-sm text-green-700">您已成功登录到Wonder系统，现在可以开始使用所有功能。</p>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="hidden lg:block">
                  <div className="w-20 h-20 bg-gradient-to-br from-blue-100 to-purple-100 rounded-2xl flex items-center justify-center">
                    <svg className="w-10 h-10 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>

            {/* 用户信息卡片 */}
            <div className="glass-card rounded-2xl p-8 card-hover">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h3 className="text-xl font-bold text-gray-900">账户信息</h3>
                  <p className="text-gray-600">您的个人账户详细信息</p>
                </div>
                <div className="icon-container w-12 h-12">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-1">
                  <label className="text-sm font-medium text-gray-500 uppercase tracking-wide">用户ID</label>
                  <div className="bg-white/70 backdrop-blur-sm p-3 rounded-xl border border-gray-200/50">
                    <p className="text-sm font-mono text-gray-900">{user?.id}</p>
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-sm font-medium text-gray-500 uppercase tracking-wide">姓名</label>
                  <div className="bg-white/70 backdrop-blur-sm p-3 rounded-xl border border-gray-200/50">
                    <p className="text-sm text-gray-900 font-medium">{user?.name}</p>
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-sm font-medium text-gray-500 uppercase tracking-wide">邮箱地址</label>
                  <div className="bg-white/70 backdrop-blur-sm p-3 rounded-xl border border-gray-200/50">
                    <p className="text-sm text-gray-900">{user?.email}</p>
                  </div>
                </div>
                <div className="space-y-1">
                  <label className="text-sm font-medium text-gray-500 uppercase tracking-wide">注册时间</label>
                  <div className="bg-white/70 backdrop-blur-sm p-3 rounded-xl border border-gray-200/50">
                    <p className="text-sm text-gray-900">
                      {user?.created_at ? new Date(user.created_at).toLocaleString('zh-CN') : ''}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* 功能区域 */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">个人资料</CardTitle>
                  <CardDescription>
                    查看和编辑您的个人信息
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Button className="w-full" variant="outline" onClick={() => setIsProfileModalOpen(true)}>
                    编辑资料
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">安全设置</CardTitle>
                  <CardDescription>
                    管理您的密码和安全选项
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Button className="w-full" variant="outline" onClick={() => setIsPasswordModalOpen(true)}>
                    安全设置
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">系统状态</CardTitle>
                  <CardDescription>
                    查看系统运行状态
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center space-x-2">
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                    <span className="text-sm text-gray-600">系统正常运行</span>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* 监控和日志管理 */}
            <MonitoringLinks />
          </div>
        </div>
      </main>

      {isProfileModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
          <div className="relative w-full max-w-lg rounded-3xl bg-white p-8 shadow-2xl">
            <button
              onClick={() => setIsProfileModalOpen(false)}
              className="absolute right-4 top-4 text-gray-400 transition hover:text-gray-600"
              aria-label="关闭"
            >
              <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 6l8 8m0-8l-8 8" />
              </svg>
            </button>

            <div className="mb-6">
              <h2 className="text-2xl font-semibold text-gray-900">编辑个人资料</h2>
              <p className="mt-1 text-sm text-gray-500">更新您的基本信息，保持资料最新。</p>
            </div>

            <form onSubmit={handleProfileSubmit} className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="profile-name">姓名</Label>
                <Input id="profile-name" placeholder="请输入姓名" {...profileForm.register('name')} />
                {profileForm.formState.errors.name && (
                  <p className="text-sm text-red-500">{profileForm.formState.errors.name.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="profile-email">邮箱地址</Label>
                <Input
                  id="profile-email"
                  type="email"
                  placeholder="请输入邮箱"
                  {...profileForm.register('email')}
                />
                {profileForm.formState.errors.email && (
                  <p className="text-sm text-red-500">{profileForm.formState.errors.email.message}</p>
                )}
              </div>

              {profileStatus && (
                <div
                  className={`rounded-xl border px-4 py-3 text-sm ${
                    profileStatus.type === 'success'
                      ? 'border-green-200 bg-green-50 text-green-700'
                      : 'border-red-200 bg-red-50 text-red-600'
                  }`}
                >
                  {profileStatus.message}
                </div>
              )}

              <div className="flex justify-end space-x-3">
                <Button type="button" variant="outline" onClick={() => setIsProfileModalOpen(false)}>
                  取消
                </Button>
                <Button type="submit" disabled={profileForm.formState.isSubmitting}>
                  {profileForm.formState.isSubmitting ? '保存中...' : '保存更改'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {isPasswordModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
          <div className="relative w-full max-w-lg rounded-3xl bg-white p-8 shadow-2xl">
            <button
              onClick={() => setIsPasswordModalOpen(false)}
              className="absolute right-4 top-4 text-gray-400 transition hover:text-gray-600"
              aria-label="关闭"
            >
              <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 6l8 8m0-8l-8 8" />
              </svg>
            </button>

            <div className="mb-6">
              <h2 className="text-2xl font-semibold text-gray-900">安全设置</h2>
              <p className="mt-1 text-sm text-gray-500">更新账户密码，提升安全性。</p>
            </div>

            <form onSubmit={handlePasswordSubmit} className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="old-password">当前密码</Label>
                <Input
                  id="old-password"
                  type="password"
                  placeholder="请输入当前密码"
                  {...passwordForm.register('oldPassword')}
                />
                {passwordForm.formState.errors.oldPassword && (
                  <p className="text-sm text-red-500">{passwordForm.formState.errors.oldPassword.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="new-password">新密码</Label>
                <Input
                  id="new-password"
                  type="password"
                  placeholder="请输入新密码"
                  {...passwordForm.register('newPassword')}
                />
                {passwordForm.formState.errors.newPassword && (
                  <p className="text-sm text-red-500">{passwordForm.formState.errors.newPassword.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirm-password">确认新密码</Label>
                <Input
                  id="confirm-password"
                  type="password"
                  placeholder="请再次输入新密码"
                  {...passwordForm.register('confirmPassword')}
                />
                {passwordForm.formState.errors.confirmPassword && (
                  <p className="text-sm text-red-500">{passwordForm.formState.errors.confirmPassword.message}</p>
                )}
              </div>

              {passwordStatus && (
                <div
                  className={`rounded-xl border px-4 py-3 text-sm ${
                    passwordStatus.type === 'success'
                      ? 'border-green-200 bg-green-50 text-green-700'
                      : 'border-red-200 bg-red-50 text-red-600'
                  }`}
                >
                  {passwordStatus.message}
                </div>
              )}

              <div className="flex justify-end space-x-3">
                <Button type="button" variant="outline" onClick={() => setIsPasswordModalOpen(false)}>
                  取消
                </Button>
                <Button type="submit" disabled={passwordForm.formState.isSubmitting}>
                  {passwordForm.formState.isSubmitting ? '更新中...' : '更新密码'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
