import React, { useState } from 'react';
import { Row, Col, Card, Form, Button, Alert } from 'react-bootstrap';
import { analyticsAPI } from '../services/api';

const Reports = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [filters, setFilters] = useState({
    startDate: '',
    endDate: '',
    reportType: 'banks'
  });

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleGenerateReport = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);

      const filterData = {
        start_date: filters.startDate,
        end_date: filters.endDate
      };

      let response;
      switch (filters.reportType) {
        case 'banks':
          response = await analyticsAPI.generateBanksReport(filterData);
          break;
        case 'categories':
          response = await analyticsAPI.getCategoriesSummary(filterData);
          break;
        case 'status':
          response = await analyticsAPI.getStatusSummary(filterData);
          break;
        case 'dynamics':
          response = await analyticsAPI.getDynamicsByPeriod('month', filterData);
          break;
        case 'compare':
          response = await analyticsAPI.getIncomeExpenseComparison(filterData);
          break;
        default:
          throw new Error('Неизвестный тип отчета');
      }

      if (filters.reportType === 'banks') {
        // Скачиваем PDF отчет
        const reportResponse = await analyticsAPI.downloadReport(response.data.report_id);
        const blob = new Blob([reportResponse.data], { type: 'application/pdf' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `report-${response.data.report_id}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      }

      setSuccess('Отчет успешно сгенерирован');
    } catch (err) {
      console.error('Ошибка при генерации отчета:', err);
      setError(err.response?.data?.message || 'Ошибка при генерации отчета');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mt-4">
      <Card>
        <Card.Header>
          <h2 className="mb-0">Аналитика</h2>
        </Card.Header>
        <Card.Body>
          <Form onSubmit={handleGenerateReport}>
            <Row>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Начальная дата</Form.Label>
                  <Form.Control
                    type="date"
                    name="startDate"
                    value={filters.startDate}
                    onChange={handleFilterChange}
                    required
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Конечная дата</Form.Label>
                  <Form.Control
                    type="date"
                    name="endDate"
                    value={filters.endDate}
                    onChange={handleFilterChange}
                    required
                  />
                </Form.Group>
              </Col>
              <Col md={4}>
                <Form.Group className="mb-3">
                  <Form.Label>Тип отчета</Form.Label>
                  <Form.Select
                    name="reportType"
                    value={filters.reportType}
                    onChange={handleFilterChange}
                  >
                    <option value="banks">По банкам</option>
                    <option value="categories">По категориям</option>
                    <option value="status">По статусам</option>
                    <option value="dynamics">Динамика</option>
                    <option value="compare">Сравнение доходов/расходов</option>
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>
            <Row>
              <Col className="text-end">
                <Button 
                  variant="primary" 
                  type="submit" 
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                      Генерация...
                    </>
                  ) : (
                    'Сгенерировать отчет'
                  )}
                </Button>
              </Col>
            </Row>
          </Form>

          {error && (
            <Alert variant="danger" className="mt-3">
              {error}
            </Alert>
          )}

          {success && (
            <Alert variant="success" className="mt-3">
              {success}
            </Alert>
          )}
        </Card.Body>
      </Card>

      <Card className="mt-4">
        <Card.Header>
          <h3 className="mb-0">Доступные отчеты</h3>
        </Card.Header>
        <Card.Body>
          <Row>
            <Col md={4}>
              <Card className="mb-3">
                <Card.Body>
                  <h5 className="card-title">Отчет по банкам</h5>
                  <p className="card-text">
                    Детальный анализ транзакций по банкам и счетам
                  </p>
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card className="mb-3">
                <Card.Body>
                  <h5 className="card-title">Отчет по категориям</h5>
                  <p className="card-text">
                    Анализ расходов и доходов по категориям
                  </p>
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card className="mb-3">
                <Card.Body>
                  <h5 className="card-title">Отчет по статусам</h5>
                  <p className="card-text">
                    Анализ транзакций по статусам
                  </p>
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card className="mb-3">
                <Card.Body>
                  <h5 className="card-title">Динамика</h5>
                  <p className="card-text">
                    Анализ динамики транзакций по периодам
                  </p>
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card className="mb-3">
                <Card.Body>
                  <h5 className="card-title">Сравнение</h5>
                  <p className="card-text">
                    Сравнение доходов и расходов
                  </p>
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </Card.Body>
      </Card>
    </div>
  );
};

export default Reports; 