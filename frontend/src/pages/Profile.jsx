import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Form, Button, Alert } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../services/api';

const Profile = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [subjectTypes, setSubjectTypes] = useState([]);
  const [profile, setProfile] = useState({
    login: '',
    user_type: ''
  });

  useEffect(() => {
    fetchSubjectTypes();
  }, []);

  const fetchSubjectTypes = async () => {
    try {
      const response = await authAPI.getSubjectTypes();
      setSubjectTypes(response.data);
    } catch (err) {
      console.error('Ошибка при загрузке типов пользователей:', err);
      setError(err.response?.data?.message || 'Ошибка при загрузке типов пользователей');
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProfile(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className="container mt-4">
      <Card>
        <Card.Header>
          <h2 className="mb-0">Профиль</h2>
        </Card.Header>
        <Card.Body>
          {error && (
            <Alert variant="danger" className="mb-4">
              {error}
            </Alert>
          )}

          {success && (
            <Alert variant="success" className="mb-4">
              {success}
            </Alert>
          )}

          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Логин</Form.Label>
              <Form.Control
                type="text"
                name="login"
                value={profile.login}
                onChange={handleInputChange}
                disabled
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>Тип пользователя</Form.Label>
              <Form.Control
                type="text"
                name="user_type"
                value={profile.user_type}
                onChange={handleInputChange}
                disabled
              />
            </Form.Group>

            <Button variant="danger" onClick={handleLogout}>
              Выйти
            </Button>
          </Form>
        </Card.Body>
      </Card>
    </div>
  );
};

export default Profile; 