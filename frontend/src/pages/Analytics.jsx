import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Form } from 'react-bootstrap';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { analyticsAPI } from '../services/api';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8'];

const Analytics = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [period, setPeriod] = useState('month');
  const [dynamicsData, setDynamicsData] = useState([]);
  const [categoryIncomeSummary, setCategoryIncomeSummary] = useState([]);
  const [categoryExpenseSummary, setCategoryExpenseSummary] = useState([]);

  useEffect(() => {
    fetchAnalyticsData();
  }, [period]);

  const fetchAnalyticsData = async () => {
    try {
      setLoading(true);
      const now = new Date();
      const startDate = new Date(now.getFullYear(), now.getMonth(), 1);
      const endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0);

      const filter = {
        date: {
          from: startDate.toISOString().split('T')[0],
          to: endDate.toISOString().split('T')[0]
        }
      };

      const [dynamicsResponse, incomeResponse, expenseResponse] = await Promise.all([
        analyticsAPI.getDynamicsByPeriod(period, filter),
        analyticsAPI.getCategoriesSummary('credit', filter),
        analyticsAPI.getCategoriesSummary('debit', filter)
      ]);

      setDynamicsData(dynamicsResponse.data?.data || []);
      setCategoryIncomeSummary(incomeResponse.data?.data || []);
      setCategoryExpenseSummary(expenseResponse.data?.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке данных:', err);
      setError(err.response?.data?.message || 'Ошибка при загрузке данных');
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

  if (error) {
    return (
      <div className="alert alert-danger" role="alert">
        {error}
      </div>
    );
  }

  return (
    <>
      <Row className="mb-4">
        <Col>
          <h1 className="h3 mb-0 text-gray-800">Аналитика</h1>
        </Col>
        <Col xs="auto">
          <Form.Select
            value={period}
            onChange={(e) => setPeriod(e.target.value)}
            style={{ width: '200px' }}
          >
            <option value="week">По неделям</option>
            <option value="month">По месяцам</option>
            <option value="quarter">По кварталам</option>
            <option value="year">По годам</option>
          </Form.Select>
        </Col>
      </Row>

      <Row className="mb-4">
        <Col lg={8}>
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Динамика доходов и расходов</Card.Title>
              <div style={{ height: '400px' }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={dynamicsData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="date" 
                      tickFormatter={date => {
                        try {
                          return new Date(date).toLocaleDateString('ru-RU', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric'
                          });
                        } catch (e) {
                          return date;
                        }
                      }} 
                    />
                    <YAxis />
                    <Tooltip 
                      labelFormatter={date => {
                        try {
                          return new Date(date).toLocaleDateString('ru-RU', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric'
                          });
                        } catch (e) {
                          return date;
                        }
                      }} 
                    />
                    <Legend />
                    <Bar dataKey="value" name="Сумма" fill="#28a745" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </Card.Body>
          </Card>
        </Col>

        <Col lg={4}>
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Распределение по категориям (доходы)</Card.Title>
              <div style={{ height: '200px' }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={categoryIncomeSummary}
                      dataKey="value"
                      nameKey="category"
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      label={(entry) => entry.category}
                    >
                      {categoryIncomeSummary.map((entry, index) => (
                        <Cell key={`cell-income-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              <Card.Title className="text-primary mb-4 mt-4">Распределение по категориям (расходы)</Card.Title>
              <div style={{ height: '200px' }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={categoryExpenseSummary}
                      dataKey="value"
                      nameKey="category"
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      label={(entry) => entry.category}
                    >
                      {categoryExpenseSummary.map((entry, index) => (
                        <Cell key={`cell-expense-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row>
        <Col md={6}>
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Топ категорий расходов</Card.Title>
              {categoryExpenseSummary
                .filter(category => category && category.value !== undefined)
                .sort((a, b) => b.value - a.value)
                .slice(0, 5)
                .map((category, index) => (
                  <div key={index} className="mb-3">
                    <div className="d-flex justify-content-between mb-1">
                      <span>{category.category || 'Без категории'}</span>
                      <span className="text-danger">
                        {(category.value || 0).toLocaleString('ru-RU', {
                          style: 'currency',
                          currency: 'RUB'
                        })}
                      </span>
                    </div>
                    <div className="progress" style={{ height: '10px' }}>
                      <div
                        className="progress-bar bg-danger"
                        role="progressbar"
                        style={{
                          width: `${((category.value || 0) / (categoryExpenseSummary.reduce((sum, cat) => sum + (cat.value || 0), 0) || 1)) * 100}%`
                        }}
                      />
                    </div>
                  </div>
                ))}
            </Card.Body>
          </Card>
        </Col>

        <Col md={6}>
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Топ категорий доходов</Card.Title>
              {categoryIncomeSummary
                .filter(category => category && category.value !== undefined)
                .sort((a, b) => b.value - a.value)
                .slice(0, 5)
                .map((category, index) => (
                  <div key={index} className="mb-3">
                    <div className="d-flex justify-content-between mb-1">
                      <span>{category.category || 'Без категории'}</span>
                      <span className="text-success">
                        {(category.value || 0).toLocaleString('ru-RU', {
                          style: 'currency',
                          currency: 'RUB'
                        })}
                      </span>
                    </div>
                    <div className="progress" style={{ height: '10px' }}>
                      <div
                        className="progress-bar bg-success"
                        role="progressbar"
                        style={{
                          width: `${((category.value || 0) / (categoryIncomeSummary.reduce((sum, cat) => sum + (cat.value || 0), 0) || 1)) * 100}%`
                        }}
                      />
                    </div>
                  </div>
                ))}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </>
  );
};

export default Analytics;