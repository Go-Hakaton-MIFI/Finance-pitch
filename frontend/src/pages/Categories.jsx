import React, { useState, useEffect } from 'react';
import { Container, Table, Alert } from 'react-bootstrap';
import { categoriesAPI } from '../services/api';

const Categories = () => {
  const [categories, setCategories] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      const response = await categoriesAPI.getAll();
      setCategories(response.data);
    } catch (err) {
      setError('Ошибка при загрузке категорий');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Container className="mt-4">
        <div className="text-center">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      </Container>
    );
  }

  return (
    <Container className="py-4">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h1>Категории</h1>
      </div>

      {error && <Alert variant="danger" className="mb-4">{error}</Alert>}

      <Table responsive striped bordered hover>
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
          </tr>
        </thead>
        <tbody>
          {categories.map((category) => (
            <tr key={category.id}>
              <td>{category.id}</td>
              <td>{category.name}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Container>
  );
};

export default Categories; 