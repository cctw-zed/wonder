'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { Button } from '@/components/ui/button';

interface MonitoringTool {
  name: string;
  url: string;
  description: string;
  icon: string;
}

const monitoringTools: MonitoringTool[] = [
  {
    name: 'Grafana',
    url: 'http://localhost:3000',
    description: '监控仪表盘',
    icon: '📊',
  },
  {
    name: 'Prometheus',
    url: 'http://localhost:9090',
    description: '指标查询',
    icon: '📈',
  },
  {
    name: 'Kibana',
    url: 'http://localhost:5601',
    description: '日志分析',
    icon: '🔍',
  },
  {
    name: 'cAdvisor',
    url: 'http://localhost:8081',
    description: '容器监控',
    icon: '🐳',
  },
];

export default function MonitoringDropdown() {
  const [isOpen, setIsOpen] = useState(false);
  const [isPortalReady, setIsPortalReady] = useState(false);
  const [dropdownPosition, setDropdownPosition] = useState<{ top: number; left: number } | null>(null);
  const triggerRef = useRef<HTMLButtonElement | null>(null);

  useEffect(() => {
    setIsPortalReady(true);
  }, []);

  const updatePosition = useCallback(() => {
    if (!triggerRef.current) return;

    const rect = triggerRef.current.getBoundingClientRect();
    const viewportPadding = 12;

    const availableWidth = window.innerWidth;
    const desiredLeft = rect.left;
    const adjustedLeft = Math.min(
      Math.max(desiredLeft, viewportPadding),
      Math.max(viewportPadding, availableWidth - viewportPadding - 256)
    );

    setDropdownPosition({
      top: rect.bottom + 8,
      left: adjustedLeft,
    });
  }, []);

  useEffect(() => {
    if (!isOpen) return;

    updatePosition();

    const handleReposition = () => updatePosition();

    window.addEventListener('resize', handleReposition);
    window.addEventListener('scroll', handleReposition, true);

    return () => {
      window.removeEventListener('resize', handleReposition);
      window.removeEventListener('scroll', handleReposition, true);
    };
  }, [isOpen, updatePosition]);

  const toggleDropdown = () => {
    if (!isOpen) {
      updatePosition();
    }
    setIsOpen((prev) => !prev);
  };

  const handleOpenTool = (url: string, name: string) => {
    try {
      window.open(url, '_blank', 'noopener,noreferrer');
      setIsOpen(false);
    } catch (error) {
      console.error(`Failed to open ${name}:`, error);
    }
  };

  return (
    <div className="relative inline-flex">
      <Button
        variant="outline"
        size="sm"
        ref={triggerRef}
        onClick={toggleDropdown}
        className="flex items-center space-x-1"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
          />
        </svg>
        <span className="hidden sm:inline">监控</span>
        <svg
          className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </Button>

      {isPortalReady && isOpen && dropdownPosition &&
        createPortal(
          <>
            {/* 背景遮罩 */}
            <div
              className="fixed inset-0 z-[180]"
              onClick={() => setIsOpen(false)}
            />

            {/* 下拉菜单 */}
            <div
              className="fixed z-[200] w-64 rounded-md border border-gray-200 bg-white shadow-lg"
              style={{ top: dropdownPosition.top, left: dropdownPosition.left }}
            >
              <div className="py-2">
                <div className="px-4 py-2 text-sm font-medium text-gray-700 border-b border-gray-100">
                  监控工具
                </div>

                {monitoringTools.map((tool) => (
                  <button
                    key={tool.name}
                    onClick={() => handleOpenTool(tool.url, tool.name)}
                    className="flex w-full items-center space-x-3 px-4 py-3 text-left transition-colors hover:bg-gray-50"
                  >
                    <span className="text-lg">{tool.icon}</span>
                    <div className="min-w-0 flex-1">
                      <div className="text-sm font-medium text-gray-900">
                        {tool.name}
                      </div>
                      <div className="text-xs text-gray-500 truncate">
                        {tool.description}
                      </div>
                    </div>
                    <svg
                      className="h-4 w-4 text-gray-400"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                      />
                    </svg>
                  </button>
                ))}

                <div className="mt-2 border-t border-gray-100 pt-2">
                  <div className="px-4 py-2 text-xs text-gray-500">
                    点击任意工具在新窗口中打开
                  </div>
                </div>
              </div>
            </div>
          </>,
          document.body
        )}
    </div>
  );
}
