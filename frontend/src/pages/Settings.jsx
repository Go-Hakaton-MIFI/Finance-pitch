import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Form, Button, Alert, Nav, Tab } from 'react-bootstrap';
import { settingsAPI } from '../services/api';

const Settings = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [settings, setSettings] = useState({
    notifications: {
      email: true,
      push: true,
      dailyReport: true,
      weeklyReport: true
    },
    currency: 'RUB',
    language: 'ru',
    theme: 'light'
  });

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    try {
      setLoading(true);
      const response = await settingsAPI.getSettings();
      setSettings(response.data);
    } catch (err) {
      console.error('Ошибка при загрузке настроек:', err);
      setError(err.response?.data?.message || 'Ошибка при загрузке настроек');
    } finally {
      setLoading(false);
    }
  };

  const handleSettingChange = (category, setting, value) => {
    setSettings(prev => ({
      ...prev,
      [category]: {
        ...prev[category],
        [setting]: value
      }
    }));
  };

  const handleSaveSettings = async () => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);

      await settingsAPI.updateSettings(settings);
      setSuccess('Настройки успешно сохранены');
    } catch (err) {
      console.error('Ошибка при сохранении настроек:', err);
      setError(err.response?.data?.message || 'Ошибка при сохранении настроек');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="text-center">
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Загрузка...</span>
        </div>
      </div>
    );
  }

  return (
    <>
      <Row className="mb-4">
        <Col>
          <h1 className="h3 mb-0 text-gray-800">Настройки</h1>
        </Col>
        <Col xs="auto">
          <Button 
            variant="primary" 
            onClick={handleSaveSettings}
            disabled={loading}
          >
            {loading ? (
              <>
                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                Сохранение...
              </>
            ) : (
              <>
                <i className="bi bi-save me-2"></i>
                Сохранить изменения
              </>
            )}
          </Button>
        </Col>
      </Row>

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

      <Card className="shadow">
        <Card.Body>
          <Tab.Container defaultActiveKey="notifications">
            <Row>
              <Col md={3}>
                <Nav variant="pills" className="flex-column">
                  <Nav.Item>
                    <Nav.Link eventKey="notifications">
                      <i className="bi bi-bell me-2"></i>
                      Уведомления
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="general">
                      <i className="bi bi-gear me-2"></i>
                      Общие
                    </Nav.Link>
                  </Nav.Item>
                  <Nav.Item>
                    <Nav.Link eventKey="appearance">
                      <i className="bi bi-palette me-2"></i>
                      Внешний вид
                    </Nav.Link>
                  </Nav.Item>
                </Nav>
              </Col>
              <Col md={9}>
                <Tab.Content>
                  <Tab.Pane eventKey="notifications">
                    <h5 className="mb-4">Настройки уведомлений</h5>
                    <Form>
                      <Form.Group className="mb-3">
                        <Form.Check
                          type="switch"
                          id="email-notifications"
                          label="Email уведомления"
                          checked={settings.notifications.email}
                          onChange={(e) => handleSettingChange('notifications', 'email', e.target.checked)}
                        />
                      </Form.Group>
                      <Form.Group className="mb-3">
                        <Form.Check
                          type="switch"
                          id="push-notifications"
                          label="Push уведомления"
                          checked={settings.notifications.push}
                          onChange={(e) => handleSettingChange('notifications', 'push', e.target.checked)}
                        />
                      </Form.Group>
                      <Form.Group className="mb-3">
                        <Form.Check
                          type="switch"
                          id="daily-report"
                          label="Ежедневный отчет"
                          checked={settings.notifications.dailyReport}
                          onChange={(e) => handleSettingChange('notifications', 'dailyReport', e.target.checked)}
                        />
                      </Form.Group>
                      <Form.Group className="mb-3">
                        <Form.Check
                          type="switch"
                          id="weekly-report"
                          label="Еженедельный отчет"
                          checked={settings.notifications.weeklyReport}
                          onChange={(e) => handleSettingChange('notifications', 'weeklyReport', e.target.checked)}
                        />
                      </Form.Group>
                    </Form>
                  </Tab.Pane>
                  <Tab.Pane eventKey="general">
                    <h5 className="mb-4">Общие настройки</h5>
                    <Form>
                      <Form.Group className="mb-3">
                        <Form.Label>Валюта</Form.Label>
                        <Form.Select
                          value={settings.currency}
                          onChange={(e) => handleSettingChange('currency', null, e.target.value)}
                        >
                          <option value="RUB">Рубль (₽)</option>
                          <option value="USD">Доллар ($)</option>
                          <option value="EUR">Евро (€)</option>
                        </Form.Select>
                      </Form.Group>
                      <Form.Group className="mb-3">
                        <Form.Label>Язык</Form.Label>
                        <Form.Select
                          value={settings.language}
                          onChange={(e) => handleSettingChange('language', null, e.target.value)}
                        >
                          <option value="ru">Русский</option>
                          <option value="en">English</option>
                        </Form.Select>
                      </Form.Group>
                    </Form>
                  </Tab.Pane>
                  <Tab.Pane eventKey="appearance">
                    <h5 className="mb-4">Настройки внешнего вида</h5>
                    <Form>
                      <Form.Group className="mb-3">
                        <Form.Label>Тема</Form.Label>
                        <Form.Select
                          value={settings.theme}
                          onChange={(e) => handleSettingChange('theme', null, e.target.value)}
                        >
                          <option value="light">Светлая</option>
                          <option value="dark">Темная</option>
                          <option value="system">Системная</option>
                        </Form.Select>
                      </Form.Group>
                    </Form>
                  </Tab.Pane>
                </Tab.Content>
              </Col>
            </Row>
          </Tab.Container>
        </Card.Body>
      </Card>
    </>
  );
};

export default Settings; 