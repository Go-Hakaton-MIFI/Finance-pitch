import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Formik, Form } from 'formik';
import * as Yup from 'yup';
import { Container, Card, Form as BootstrapForm, Button, Alert } from 'react-bootstrap';
import { authAPI } from '../services/api';

const registerSchema = Yup.object().shape({
  loginName: Yup.string()
    .min(3, 'Логин должен быть не менее 3 символов')
    .max(50, 'Логин должен быть не более 50 символов')
    .required('Логин обязателен'),
  password: Yup.string()
    .min(6, 'Пароль должен быть не менее 6 символов')
    .required('Пароль обязателен'),
  userType: Yup.string()
    .required('Тип пользователя обязателен')
    .oneOf(['ФЛ', 'ЮЛ'], 'Тип пользователя должен быть выбран'),
  partName: Yup.string()
    .required('Имя обязательно')
    .min(3, 'Имя должно быть не менее 3 символов')
    .max(50, 'Имя должно быть не более 50 символов'),
  inn: Yup.string()
    .required('ИНН обязателен')
    .matches(/^\d{11}$/, 'ИНН должен содержать 11 цифр'),
  phone: Yup.string()
    .required('Телефон обязателен')
    .matches(/^\+7\d{10}$/, 'Телефон должен быть в формате +7XXXXXXXXXX'),
  bank: Yup.string()
    .required('Название банка обязательно'),
  account: Yup.string()
    .required('Номер счета обязателен')
    .matches(/^\d{20}$/, 'Номер счета должен содержать 20 цифр')
});

const Register = () => {
  const navigate = useNavigate();
  const [error, setError] = useState('');
  const [subjectTypes, setSubjectTypes] = useState([]);

  useEffect(() => {
    const fetchSubjectTypes = async () => {
      try {
        const response = await authAPI.getSubjectTypes();
        setSubjectTypes(response.data || []);
      } catch (err) {
        console.error('Ошибка при получении типов пользователей:', err);
        setError('Ошибка при загрузке типов пользователей');
      }
    };
    fetchSubjectTypes();
  }, []);

  const handleSubmit = async (values, { setSubmitting }) => {
    try {
      await authAPI.register(values);
      navigate('/login');
    } catch (err) {
      setError(err.response?.data?.message || 'Ошибка при регистрации');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Container className="d-flex justify-content-center align-items-center" style={{ minHeight: '100vh' }}>
      <Card style={{ width: '500px' }}>
        <Card.Body className="p-4">
          <Card.Title className="text-center mb-4">Регистрация</Card.Title>
          {error && <Alert variant="danger">{error}</Alert>}
          <Formik
            initialValues={{
              loginName: '',
              password: '',
              userType: '',
              partName: '',
              inn: '',
              phone: '',
              bank: '',
              account: ''
            }}
            validationSchema={registerSchema}
            onSubmit={handleSubmit}
          >
            {({ handleChange, handleBlur, values, touched, errors, isSubmitting }) => (
              <Form>
                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Тип пользователя</BootstrapForm.Label>
                  <BootstrapForm.Select
                    name="userType"
                    value={values.userType}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.userType && errors.userType}
                  >
                    <option value="">Выберите тип пользователя</option>
                    {subjectTypes.map((type) => (
                      <option key={type.subjectType} value={type.subjectType === 'INDIVIDUAL' ? 'ФЛ' : 'ЮЛ'}>
                        {type.subjectName}
                      </option>
                    ))}
                  </BootstrapForm.Select>
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.userType}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Логин</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="loginName"
                    value={values.loginName}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.loginName && errors.loginName}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.loginName}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Пароль</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="password"
                    name="password"
                    value={values.password}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.password && errors.password}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.password}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Имя</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="partName"
                    value={values.partName}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.partName && errors.partName}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.partName}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>ИНН</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="inn"
                    value={values.inn}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.inn && errors.inn}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.inn}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Телефон</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="phone"
                    value={values.phone}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.phone && errors.phone}
                    placeholder="+7XXXXXXXXXX"
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.phone}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Банк</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="bank"
                    value={values.bank}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.bank && errors.bank}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.bank}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <BootstrapForm.Group className="mb-3">
                  <BootstrapForm.Label>Номер счета</BootstrapForm.Label>
                  <BootstrapForm.Control
                    type="text"
                    name="account"
                    value={values.account}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    isInvalid={touched.account && errors.account}
                  />
                  <BootstrapForm.Control.Feedback type="invalid">
                    {errors.account}
                  </BootstrapForm.Control.Feedback>
                </BootstrapForm.Group>

                <Button
                  variant="primary"
                  type="submit"
                  className="w-100"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'Регистрация...' : 'Зарегистрироваться'}
                </Button>
              </Form>
            )}
          </Formik>
        </Card.Body>
      </Card>
    </Container>
  );
};

export default Register;