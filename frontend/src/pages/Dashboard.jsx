import React, { useEffect, useState } from 'react';
import { Row, Col, Card, Alert } from 'react-bootstrap';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { analyticsAPI } from '../services/api';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8'];

const formatDate = (date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [dynamicsData, setDynamicsData] = useState([]);
  const [categorySummary, setCategorySummary] = useState([]);
  const [comparison, setComparison] = useState({
    income: 0,
    expense: 0,
    balance: 0
  });

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        setError(null);

        const now = new Date();
        const startDate = new Date(now.getFullYear(), now.getMonth(), 1);
        const endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0);

        const filter = {
          date: {
            from: startDate.toISOString(),
            to: endDate.toISOString()
          }
        };

        console.log('Отправляем запросы с фильтром:', filter);

        try {
          const [dynamicsResponse, incomeResponse, expenseResponse] = await Promise.all([
            analyticsAPI.getDynamicsByPeriod('month', filter),
            analyticsAPI.getCategoriesSummary('credit', filter),
            analyticsAPI.getCategoriesSummary('debit', filter)
          ]);

          console.log('Dynamics response:', dynamicsResponse.data);
          console.log('Income response:', incomeResponse.data);
          console.log('Expense response:', expenseResponse.data);

          setDynamicsData(dynamicsResponse.data?.data || []);
          
          const incomeData = (incomeResponse.data?.data || []).map(item => ({
            ...item,
            type: 'INCOME'
          }));
          const expenseData = (expenseResponse.data?.data || []).map(item => ({
            ...item,
            type: 'EXPENSE'
          }));
          
          setCategorySummary([...incomeData, ...expenseData]);
          
          const totalIncome = incomeData.reduce((sum, item) => sum + (item.value || 0), 0);
          const totalExpense = expenseData.reduce((sum, item) => sum + (item.value || 0), 0);
          
          setComparison({
            income: totalIncome,
            expense: totalExpense,
            balance: totalIncome - totalExpense
          });
        } catch (err) {
          console.error('Ошибка при загрузке данных:', err);
          if (err.response?.data?.error) {
            setError(err.response.data.error);
          } else if (err.response?.data?.message) {
            setError(err.response.data.message);
          } else {
            setError('Ошибка при загрузке данных');
          }
        }
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

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
      <Alert variant="danger">
        <Alert.Heading>Ошибка при загрузке данных</Alert.Heading>
        <p>{error}</p>
      </Alert>
    );
  }

  return (
    <div className="container mt-4">
      <Row className="mb-4">
        <Col lg={4} md={6} className="mb-4">
          <Card className="h-100">
            <Card.Body>
              <Card.Title className="text-primary">Доходы за месяц</Card.Title>
              <h3 className="text-success">
                {comparison.income.toLocaleString('ru-RU', {
                  style: 'currency',
                  currency: 'RUB'
                })}
              </h3>
              {comparison.income === 0 && (
                <small className="text-muted">Данные временно недоступны</small>
              )}
            </Card.Body>
          </Card>
        </Col>

        <Col lg={4} md={6} className="mb-4">
          <Card className="h-100">
            <Card.Body>
              <Card.Title className="text-primary">Расходы за месяц</Card.Title>
              <h3 className="text-danger">
                {comparison.expense.toLocaleString('ru-RU', {
                  style: 'currency',
                  currency: 'RUB'
                })}
              </h3>
              {comparison.expense === 0 && (
                <small className="text-muted">Данные временно недоступны</small>
              )}
            </Card.Body>
          </Card>
        </Col>

        <Col lg={4} md={6} className="mb-4">
          <Card className="h-100">
            <Card.Body>
              <Card.Title className="text-primary">Баланс</Card.Title>
              <h3 className={comparison.balance >= 0 ? 'text-success' : 'text-danger'}>
                {comparison.balance.toLocaleString('ru-RU', {
                  style: 'currency',
                  currency: 'RUB'
                })}
              </h3>
              {comparison.balance === 0 && (
                <small className="text-muted">Данные временно недоступны</small>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row className="mb-4">
        <Col lg={8} className="mb-4">
          <Card>
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
          <Card>
            <Card.Body>
              <Card.Title className="text-primary mb-4">Распределение по категориям</Card.Title>
              <div style={{ height: '400px' }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={categorySummary}
                      dataKey="value"
                      nameKey="category"
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      label={(entry) => entry.category}
                    >
                      {categorySummary.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
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
        <Col md={6} className="mb-4">
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Топ категорий расходов</Card.Title>
              {categorySummary
                .filter(cat => cat.type === 'EXPENSE')
                .slice(0, 5)
                .map((category, index) => (
                  <div key={index} className="mb-3">
                    <div className="d-flex justify-content-between mb-1">
                      <span>{category.category}</span>
                      <span className="text-danger">
                        {category.value.toLocaleString('ru-RU', {
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
                          width: `${(category.value / comparison.expense) * 100}%`
                        }}
                      />
                    </div>
                  </div>
                ))}
            </Card.Body>
          </Card>
        </Col>

        <Col md={6} className="mb-4">
          <Card className="shadow">
            <Card.Body>
              <Card.Title className="text-primary mb-4">Топ категорий доходов</Card.Title>
              {categorySummary
                .filter(cat => cat.type === 'INCOME')
                .slice(0, 5)
                .map((category, index) => (
                  <div key={index} className="mb-3">
                    <div className="d-flex justify-content-between mb-1">
                      <span>{category.category}</span>
                      <span className="text-success">
                        {category.value.toLocaleString('ru-RU', {
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
                          width: `${(category.value / comparison.income) * 100}%`
                        }}
                      />
                    </div>
                  </div>
                ))}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;