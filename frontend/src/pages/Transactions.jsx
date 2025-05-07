import React, { useState, useEffect } from 'react';
import { Table, Button, Form, Row, Col, Card, Modal } from 'react-bootstrap';
import { transactionsAPI, categoriesAPI, statusesAPI } from '../services/api';

const Transactions = () => {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [categories, setCategories] = useState([]);
  const [statuses, setStatuses] = useState([]);
  const [filters, setFilters] = useState({
    user_type: '',
    trans_type: '',
    sender_bank: '',
    receiver_inn: '',
    receiver_phone: '',
    date_from: '',
    date_to: '',
    category_id: '',
    status_id: ''
  });
  const [newTransaction, setNewTransaction] = useState({
    user_type: 'individual',
    date_time: new Date().toISOString().split('T')[0],
    trans_type: 'debit',
    amount: '',
    category_id: '',
    status_id: '',
    sender_bank: '',
    receiver_inn: '',
    receiver_phone: '',
    comment: ''
  });

  useEffect(() => {
    fetchTransactions();
    fetchCategories();
    fetchStatuses();
  }, [filters]);

  const fetchCategories = async () => {
    try {
      const response = await categoriesAPI.getAll();
      setCategories(response.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке категорий:', err);
      setError('Ошибка при загрузке категорий');
    }
  };

  const fetchStatuses = async () => {
    try {
      const response = await statusesAPI.getAll();
      setStatuses(response.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке статусов:', err);
      setError('Ошибка при загрузке статусов');
    }
  };

  const fetchTransactions = async () => {
    try {
      setLoading(true);
      // Преобразуем фильтры в формат, ожидаемый бэкендом
      const formattedFilters = {
        user_type: filters.user_type || '',
        trans_type: filters.trans_type || '',
        sender_bank: filters.sender_bank || '',
        receiver_inn: filters.receiver_inn || '',
        receiver_phone: filters.receiver_phone || '',
        date_from: filters.date_from ? new Date(filters.date_from + 'T00:00:00Z') : null,
        date_to: filters.date_to ? new Date(filters.date_to + 'T23:59:59Z') : null,
        category_id: filters.category_id ? parseInt(filters.category_id, 10) : 0,
        status_id: filters.status_id ? parseInt(filters.status_id, 10) : 0
      };

      const response = await transactionsAPI.getAll(formattedFilters);
      setTransactions(response.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке транзакций:', err);
      setError(err.response?.data?.message || 'Ошибка при загрузке транзакций');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTransaction = async (e) => {
    e.preventDefault();
    try {
      const transactionData = {
        ...newTransaction,
        date_time: new Date(newTransaction.date_time + 'T00:00:00Z').toISOString(),
        amount: parseFloat(newTransaction.amount),
        category_id: parseInt(newTransaction.category_id, 10),
        status_id: parseInt(newTransaction.status_id, 10)
      };

      const response = await transactionsAPI.create(transactionData);
      setTransactions(prev => [response.data, ...prev]);
      setShowModal(false);
      setNewTransaction({
        user_type: 'individual',
        date_time: new Date().toISOString().split('T')[0],
        trans_type: 'debit',
        amount: '',
        category_id: '',
        status_id: '',
        sender_bank: '',
        receiver_inn: '',
        receiver_phone: '',
        comment: ''
      });
    } catch (err) {
      console.error('Ошибка при создании транзакции:', err);
      setError(err.response?.data?.error || 'Ошибка при создании транзакции');
    }
  };

  const handleNewTransactionChange = (e) => {
    const { name, value } = e.target;
    setNewTransaction(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleDeleteTransaction = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить эту транзакцию?')) return;
    try {
      await transactionsAPI.delete(id);
      setTransactions(prev => prev.filter(t => t.id !== id));
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка при удалении транзакции');
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
    <div className="container mt-4">
      <Card>
        <Card.Header className="d-flex justify-content-between align-items-center">
          <h2 className="mb-0">Транзакции</h2>
          <Button variant="primary" onClick={() => setShowModal(true)}>
            Добавить транзакцию
          </Button>
        </Card.Header>
        <Card.Body>
          {error && (
            <div className="alert alert-danger" role="alert">
              {error}
            </div>
          )}

          <Card className="mb-4">
            <Card.Body>
              <h5 className="card-title mb-3">Фильтры</h5>
              <Row>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Тип пользователя</Form.Label>
                    <Form.Select
                      name="user_type"
                      value={filters.user_type}
                      onChange={handleFilterChange}
                    >
                      <option value="">Все</option>
                      <option value="individual">Физическое лицо</option>
                      <option value="legal">Юридическое лицо</option>
                    </Form.Select>
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Тип транзакции</Form.Label>
                    <Form.Select
                      name="trans_type"
                      value={filters.trans_type}
                      onChange={handleFilterChange}
                    >
                      <option value="">Все</option>
                      <option value="credit">Доход</option>
                      <option value="debit">Расход</option>
                    </Form.Select>
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Банк отправителя</Form.Label>
                    <Form.Control
                      type="text"
                      name="sender_bank"
                      value={filters.sender_bank}
                      onChange={handleFilterChange}
                      placeholder="Введите название банка"
                    />
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>ИНН получателя</Form.Label>
                    <Form.Control
                      type="text"
                      name="receiver_inn"
                      value={filters.receiver_inn}
                      onChange={handleFilterChange}
                      placeholder="Введите ИНН"
                    />
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Телефон получателя</Form.Label>
                    <Form.Control
                      type="text"
                      name="receiver_phone"
                      value={filters.receiver_phone}
                      onChange={handleFilterChange}
                      placeholder="Введите телефон"
                    />
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Дата с</Form.Label>
                    <Form.Control
                      type="date"
                      name="date_from"
                      value={filters.date_from}
                      onChange={handleFilterChange}
                    />
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Дата по</Form.Label>
                    <Form.Control
                      type="date"
                      name="date_to"
                      value={filters.date_to}
                      onChange={handleFilterChange}
                    />
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Категория</Form.Label>
                    <Form.Select
                      name="category_id"
                      value={filters.category_id}
                      onChange={handleFilterChange}
                    >
                      <option value="">Все</option>
                      {categories.map(category => (
                        <option key={category.id} value={category.id}>
                          {category.name}
                        </option>
                      ))}
                    </Form.Select>
                  </Form.Group>
                </Col>
                <Col md={3}>
                  <Form.Group className="mb-3">
                    <Form.Label>Статус</Form.Label>
                    <Form.Select
                      name="status_id"
                      value={filters.status_id}
                      onChange={handleFilterChange}
                    >
                      <option value="">Все</option>
                      {statuses.map(status => (
                        <option key={status.id} value={status.id}>
                          {status.name}
                        </option>
                      ))}
                    </Form.Select>
                  </Form.Group>
                </Col>
              </Row>
            </Card.Body>
          </Card>

          <Table striped bordered hover responsive>
            <thead>
              <tr>
                <th>Дата</th>
                <th>Тип пользователя</th>
                <th>Тип транзакции</th>
                <th>Сумма</th>
                <th>Категория</th>
                <th>Статус</th>
                <th>Банк отправителя</th>
                <th>ИНН получателя</th>
                <th>Телефон получателя</th>
                <th>Комментарий</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {transactions.map((transaction) => (
                <tr key={transaction.id}>
                  <td>{new Date(transaction.date_time).toLocaleDateString('ru-RU', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}</td>
                  <td>{transaction.user_type === 'individual' ? 'Физ. лицо' : 'Юр. лицо'}</td>
                  <td>{transaction.trans_type === 'credit' ? 'Доход' : 'Расход'}</td>
                  <td>{transaction.amount.toLocaleString('ru-RU', {
                    style: 'currency',
                    currency: 'RUB'
                  })}</td>
                  <td>{transaction.category_name || 'Без категории'}</td>
                  <td>{transaction.status_name || 'Новый'}</td>
                  <td>{transaction.sender_bank || '-'}</td>
                  <td>{transaction.receiver_inn || '-'}</td>
                  <td>{transaction.receiver_phone || '-'}</td>
                  <td>{transaction.comment || '-'}</td>
                  <td>
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => handleDeleteTransaction(transaction.id)}
                    >
                      Удалить
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </Card.Body>
      </Card>

      <Modal show={showModal} onHide={() => setShowModal(false)} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Новая транзакция</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form onSubmit={handleCreateTransaction}>
            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Тип пользователя</Form.Label>
                  <Form.Select
                    name="user_type"
                    value={newTransaction.user_type}
                    onChange={handleNewTransactionChange}
                    required
                  >
                    <option value="individual">Физическое лицо</option>
                    <option value="legal">Юридическое лицо</option>
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Тип транзакции</Form.Label>
                  <Form.Select
                    name="trans_type"
                    value={newTransaction.trans_type}
                    onChange={handleNewTransactionChange}
                    required
                  >
                    <option value="debit">Расход</option>
                    <option value="credit">Доход</option>
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>

            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Дата</Form.Label>
                  <Form.Control
                    type="date"
                    name="date_time"
                    value={newTransaction.date_time}
                    onChange={handleNewTransactionChange}
                    required
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Сумма</Form.Label>
                  <Form.Control
                    type="number"
                    name="amount"
                    value={newTransaction.amount}
                    onChange={handleNewTransactionChange}
                    required
                    step="0.01"
                  />
                </Form.Group>
              </Col>
            </Row>

            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Категория</Form.Label>
                  <Form.Select
                    name="category_id"
                    value={newTransaction.category_id}
                    onChange={handleNewTransactionChange}
                    required
                  >
                    <option value="">Выберите категорию</option>
                    {categories.map((category) => (
                      <option key={category.id} value={category.id}>
                        {category.name}
                      </option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Статус</Form.Label>
                  <Form.Select
                    name="status_id"
                    value={newTransaction.status_id}
                    onChange={handleNewTransactionChange}
                    required
                  >
                    <option value="">Выберите статус</option>
                    {statuses.map((status) => (
                      <option key={status.id} value={status.id}>
                        {status.name}
                      </option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
            </Row>

            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Банк отправителя</Form.Label>
                  <Form.Control
                    type="text"
                    name="sender_bank"
                    value={newTransaction.sender_bank}
                    onChange={handleNewTransactionChange}
                    required
                  />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>ИНН получателя</Form.Label>
                  <Form.Control
                    type="text"
                    name="receiver_inn"
                    value={newTransaction.receiver_inn}
                    onChange={handleNewTransactionChange}
                    required
                    pattern="[0-9]{10}|[0-9]{12}"
                    title="ИНН должен содержать 10 цифр для юр. лиц или 12 цифр для физ. лиц"
                  />
                </Form.Group>
              </Col>
            </Row>

            <Row>
              <Col md={6}>
                <Form.Group className="mb-3">
                  <Form.Label>Телефон получателя</Form.Label>
                  <Form.Control
                    type="tel"
                    name="receiver_phone"
                    value={newTransaction.receiver_phone}
                    onChange={handleNewTransactionChange}
                    required
                    pattern="[0-9]{11}"
                    title="Телефон должен содержать 11 цифр"
                  />
                </Form.Group>
              </Col>
            </Row>

            <Form.Group className="mb-3">
              <Form.Label>Комментарий</Form.Label>
              <Form.Control
                as="textarea"
                name="comment"
                value={newTransaction.comment}
                onChange={handleNewTransactionChange}
                rows={3}
              />
            </Form.Group>

            <Button variant="primary" type="submit">
              Создать
            </Button>
          </Form>
        </Modal.Body>
      </Modal>
    </div>
  );
};

export default Transactions;