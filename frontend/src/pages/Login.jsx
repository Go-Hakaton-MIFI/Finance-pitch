import React, { useState } from 'react';
import { Card, Form, Button, Alert, Row, Col } from 'react-bootstrap';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { loginSuccess } from '../store/authSlice';
import { authAPI } from '../services/api';

const Login = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [formData, setFormData] = useState({
    loginName: '',
    password: ''
  });

  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);

      const response = await authAPI.login(formData);
      const { token, user } = response.data;
      
      localStorage.setItem('token', token);
      dispatch(loginSuccess({ token, user }));

      // Перенаправляем на предыдущую страницу или на главную
      const from = location.state?.from?.pathname || '/';
      navigate(from, { replace: true });
    } catch (err) {
      console.error('Ошибка при входе:', err);
      setError(err.response?.data?.message || 'Ошибка при входе в систему');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-vh-100 d-flex align-items-center justify-content-center bg-light">
      <Row className="w-100 justify-content-center">
        <Col xs={12} sm={8} md={6} lg={4}>
          <Card className="shadow">
            <Card.Body className="p-5">
              <div className="text-center mb-4">
                <i className="bi bi-wallet2 text-primary" style={{ fontSize: '3rem' }}></i>
                <h1 className="h4 mt-3">Вход в систему</h1>
                <p className="text-muted">Введите свои учетные данные</p>
              </div>

              {error && (
                <Alert variant="danger" className="mb-4">
                  {error}
                </Alert>
              )}

              <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3">
                  <Form.Label>Логин</Form.Label>
                  <Form.Control
                    type="text"
                    name="loginName"
                    value={formData.loginName}
                    onChange={handleInputChange}
                    required
                    placeholder="Введите логин"
                  />
                </Form.Group>

                <Form.Group className="mb-4">
                  <Form.Label>Пароль</Form.Label>
                  <Form.Control
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleInputChange}
                    required
                    placeholder="Введите пароль"
                  />
                </Form.Group>

                <Button
                  variant="primary"
                  type="submit"
                  className="w-100 mb-3"
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                      Вход...
                    </>
                  ) : (
                    'Войти'
                  )}
                </Button>

                <div className="text-center">
                  <p className="mb-0">
                    Нет аккаунта?{' '}
                    <Link to="/register" className="text-primary">
                      Зарегистрироваться
                    </Link>
                  </p>
                </div>
              </Form>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Login; 