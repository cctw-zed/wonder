'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface MonitoringTool {
  name: string;
  description: string;
  url: string;
  icon: React.ReactNode;
  iconColor: string;
  credentials?: string;
  details: string;
}

const monitoringTools: MonitoringTool[] = [
  {
    name: 'Grafana',
    description: '监控仪表盘',
    url: 'http://localhost:3000',
    iconColor: 'text-orange-500',
    credentials: '用户名: admin / 密码: admin',
    details: '主要监控仪表盘，包含业务指标和系统性能图表',
    icon: (
      <svg className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
      </svg>
    ),
  },
  {
    name: 'Prometheus',
    description: '指标收集',
    url: 'http://localhost:9090',
    iconColor: 'text-red-500',
    details: '指标数据源，可以查询自定义指标',
    icon: (
      <svg className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
      </svg>
    ),
  },
  {
    name: 'Kibana',
    description: '日志分析',
    url: 'http://localhost:5601',
    iconColor: 'text-purple-500',
    details: '日志查询和分析，支持全文搜索',
    icon: (
      <svg className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M4 6h16v2H4zm0 5h16v2H4zm0 5h16v2H4z"/>
      </svg>
    ),
  },
  {
    name: 'cAdvisor',
    description: '容器监控',
    url: 'http://localhost:8081',
    iconColor: 'text-blue-500',
    details: 'Docker容器资源使用情况监控',
    icon: (
      <svg className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M20,8H4V6H20M20,18H4V12H20M20,4H4C2.89,4 2,4.89 2,6V18A2,2 0 0,0 4,20H20A2,2 0 0,0 22,18V6C22,4.89 21.1,4 20,4Z"/>
      </svg>
    ),
  },
];

const getActionIcon = (toolName: string) => {
  const icons = {
    Grafana: (
      <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
      </svg>
    ),
    Prometheus: (
      <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
    ),
    Kibana: (
      <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    ),
    cAdvisor: (
      <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
    ),
  };
  return icons[toolName as keyof typeof icons] || icons.Grafana;
};

const getActionText = (toolName: string) => {
  const texts = {
    Grafana: '打开仪表盘',
    Prometheus: '查看指标',
    Kibana: '分析日志',
    cAdvisor: '容器状态',
  };
  return texts[toolName as keyof typeof texts] || '访问';
};

interface MonitoringLinksProps {
  className?: string;
}

export default function MonitoringLinks({ className = '' }: MonitoringLinksProps) {
  const handleOpenTool = (url: string, name: string) => {
    try {
      window.open(url, '_blank', 'noopener,noreferrer');
    } catch (error) {
      console.error(`Failed to open ${name}:`, error);
      // 可以在这里添加错误提示
      alert(`无法打开 ${name}，请确认服务正在运行并检查URL: ${url}`);
    }
  };

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>监控与日志管理</CardTitle>
        <CardDescription>
          访问系统监控仪表盘和日志分析工具，实时了解系统运行状态
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {monitoringTools.map((tool) => (
            <div
              key={tool.name}
              className="border rounded-lg p-4 hover:bg-gray-50 transition-all duration-200 hover:shadow-md"
            >
              <div className="flex items-center space-x-3 mb-3">
                <div className={`flex-shrink-0 ${tool.iconColor}`}>
                  {tool.icon}
                </div>
                <div className="min-w-0 flex-1">
                  <h3 className="font-medium text-gray-900 truncate">{tool.name}</h3>
                  <p className="text-sm text-gray-500 truncate">{tool.description}</p>
                </div>
              </div>

              <div className="space-y-2">
                <Button
                  className="w-full"
                  variant="outline"
                  size="sm"
                  onClick={() => handleOpenTool(tool.url, tool.name)}
                >
                  {getActionIcon(tool.name)}
                  {getActionText(tool.name)}
                </Button>

                {tool.credentials && (
                  <p className="text-xs text-gray-400 text-center">
                    {tool.credentials}
                  </p>
                )}

                <p className="text-xs text-gray-500 text-center leading-relaxed">
                  {tool.details}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* 监控系统使用说明 */}
        <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-md">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-blue-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">
                监控系统使用指南
              </h3>
              <div className="mt-2 text-sm text-blue-700">
                <div className="space-y-2">
                  <div>
                    <strong>Grafana</strong>: 点击可访问主监控仪表盘，查看系统整体运行状况、API性能指标、用户活动统计等。
                  </div>
                  <div>
                    <strong>Prometheus</strong>: 原始指标查询界面，适合开发人员进行详细的性能分析和自定义查询。
                  </div>
                  <div>
                    <strong>Kibana</strong>: 日志搜索和分析平台，可以搜索应用日志、错误信息和用户行为日志。
                  </div>
                  <div>
                    <strong>cAdvisor</strong>: 容器资源使用情况，包括CPU、内存、网络和磁盘使用率。
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 快速状态检查 */}
        <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4">
          {monitoringTools.map((tool) => (
            <div key={`status-${tool.name}`} className="text-center">
              <div className="flex items-center justify-center space-x-2">
                <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                <span className="text-sm text-gray-600">{tool.name}</span>
              </div>
              <p className="text-xs text-gray-400 mt-1">运行中</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}