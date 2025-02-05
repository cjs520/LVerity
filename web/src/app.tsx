import { ProLayout } from '@ant-design/pro-components';
import { Link, Route, Routes, useLocation, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import { Suspense, lazy } from 'react';
import { routes } from './routes';

const LoginPage = lazy(() => import('./pages/Login'));

// 路由守卫组件
const PrivateRoute = ({ element }: { element: React.ReactNode }) => {
  const token = localStorage.getItem('token');
  return token ? element : <Navigate to="/login" replace />;
};

export default function App() {
  const location = useLocation();
  const isLoginPage = location.pathname === '/login';

  // 如果是登录页面，直接渲染登录组件
  if (isLoginPage) {
    return (
      <Suspense fallback={<div style={{ padding: 24, textAlign: 'center' }}><Spin size="large" /></div>}>
        <LoginPage />
      </Suspense>
    );
  }

  return (
    <ProLayout
      title="LVerity"
      location={{
        pathname: location.pathname,
      }}
      route={{
        routes,
      }}
      menuItemRender={(item, dom) => (
        <Link to={item.path ?? '/'}>{dom}</Link>
      )}
    >
      <Suspense fallback={<div style={{ padding: 24, textAlign: 'center' }}><Spin size="large" /></div>}>
        <Routes>
          {/* 登录路由 */}
          <Route path="/login" element={<LoginPage />} />
          {/* 受保护的路由 */}
          {routes.map((route) => (
            <Route
              key={route.path}
              path={route.path}
              element={<PrivateRoute element={route.element} />}
            />
          ))}
        </Routes>
      </Suspense>
    </ProLayout>
  );
}
