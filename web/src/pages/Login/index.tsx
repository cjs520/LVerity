import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { UserOutlined, LockOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { authService } from '@/services/auth';
import styles from './index.module.less';

const LoginPage: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [captchaImage, setCaptchaImage] = useState('');
  const [captchaId, setCaptchaId] = useState('');

  // 获取验证码
  const fetchCaptcha = async () => {
    try {
      const response = await authService.getCaptcha();
      if (response.success) {
        setCaptchaImage(response.data.captcha_image);
        setCaptchaId(response.data.captcha_id);
        // 清空验证码输入框
        form.setFieldValue('captcha', '');
      } else {
        message.error('获取验证码失败，请重试');
      }
    } catch (error: any) {
      message.error('获取验证码失败，请检查网络连接');
    }
  };

  // 组件加载时获取验证码
  React.useEffect(() => {
    fetchCaptcha();
  }, []);

  // 处理登录
  const handleSubmit = async (values: any) => {
    if (!captchaId) {
      message.error('验证码已过期，请重新获取');
      return;
    }

    try {
      setLoading(true);
      
      const response = await authService.login({
        username: values.username,
        password: values.password,
        captcha: values.captcha,
        captcha_id: captchaId
      });
    
      if (response.data) {
        // 存储token和用户信息
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        message.success('登录成功');
        navigate('/');
      }
    } catch (error: any) {
      // 显示具体的错误信息
      const errorMessage = error.response?.data?.error_message || error.message || '登录失败，请重试';
      message.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.title}>LVerity 授权管理系统</div>
        <Form
          form={form}
          onFinish={handleSubmit}
          className={styles.loginForm}
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
              size="large"
              autoComplete="username"
            />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              size="large"
              autoComplete="current-password"
            />
          </Form.Item>
          <Form.Item
            name="captcha"
            rules={[{ required: true, message: '请输入验证码' }]}
          >
            <div className={styles.captchaContainer}>
              <Input
                prefix={<SafetyCertificateOutlined />}
                placeholder="验证码"
                size="large"
                maxLength={6}
                type="number"
                onKeyPress={(e) => {
                  const charCode = String.fromCharCode(e.which);
                  if (!/\d/.test(charCode)) {
                    e.preventDefault();
                  }
                }}
              />
              {captchaImage ? (
                <img
                  src={captchaImage}
                  alt="验证码"
                  className={styles.captchaImage}
                  onClick={fetchCaptcha}
                  title="点击刷新验证码"
                />
              ) : (
                <Button
                  onClick={fetchCaptcha}
                  className={styles.captchaButton}
                >
                  获取验证码
                </Button>
              )}
            </div>
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              size="large"
              loading={loading}
              block
            >
              登录
            </Button>
          </Form.Item>
        </Form>
      </div>
    </div>
  );
};

export default LoginPage;