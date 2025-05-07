import React from 'react';
import TransactionForm from '../components/TransactionForm';
import { useNavigate } from 'react-router-dom';

const AddTransaction = () => {
  const navigate = useNavigate();

  const handleSuccess = () => {
    navigate('/transactions');
  };

  return <TransactionForm onSuccess={handleSuccess} />;
};

export default AddTransaction; 