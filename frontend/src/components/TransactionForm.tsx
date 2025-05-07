import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

interface Category {
  id: number;
  name: string;
}

interface Status {
  id: number;
  name: string;
}

interface Bank {
  id: number;
  name: string;
}

interface TransactionFormData {
  user_type: 'individual' | 'legal';
  date: string;
  trans_type: 'credit' | 'debit';
  amount: number;
  category_id: number;
  status_id: number;
  sender_bank: string;
  receiver_inn: string;
  receiver_phone: string;
  comment: string;
}

interface TransactionFormProps {
  onSuccess?: () => void;
}

const TransactionForm: React.FC<TransactionFormProps> = ({ onSuccess }) => {
  const [formData, setFormData] = useState<TransactionFormData>({
    user_type: 'individual',
    date: new Date().toISOString().split('T')[0],
    trans_type: 'debit',
    amount: 0,
    category_id: 0,
    status_id: 1,
    sender_bank: '',
    receiver_inn: '',
    receiver_phone: '',
    comment: ''
  });
  const [categories, setCategories] = useState<Category[]>([]);
  const [statuses, setStatuses] = useState<Status[]>([]);
  const [banks, setBanks] = useState<Bank[]>([]);
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { token } = useAuth();

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Загрузка категорий
        const categoriesResponse = await fetch('/api/categories', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        if (!categoriesResponse.ok) {
          throw new Error('Failed to fetch categories');
        }
        const categoriesData = await categoriesResponse.json();
        setCategories(categoriesData);
        if (categoriesData.length > 0) {
          setFormData(prev => ({ ...prev, category_id: categoriesData[0].id }));
        }

        // Загрузка статусов
        const statusesResponse = await fetch('/api/statuses', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        if (!statusesResponse.ok) {
          throw new Error('Failed to fetch statuses');
        }
        const statusesData = await statusesResponse.json();
        setStatuses(statusesData);

        // Загрузка банков
        const banksResponse = await fetch('/api/banks', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        if (!banksResponse.ok) {
          throw new Error('Failed to fetch banks');
        }
        const banksData = await banksResponse.json();
        setBanks(banksData);
      } catch (err) {
        setError('Ошибка при загрузке данных');
      }
    };

    fetchData();
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/transactions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error('Ошибка при создании транзакции');
      }

      // Очистка формы после успешного создания
      setFormData({
        user_type: 'individual',
        date: new Date().toISOString().split('T')[0],
        trans_type: 'debit',
        amount: 0,
        category_id: categories[0]?.id || 0,
        status_id: 1,
        sender_bank: '',
        receiver_inn: '',
        receiver_phone: '',
        comment: ''
      });

      // Вызываем onSuccess если он передан
      if (onSuccess) {
        onSuccess();
      } else {
        // Иначе используем стандартную навигацию
        navigate('/transactions');
      }
    } catch (err) {
      setError('Ошибка при создании транзакции');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container py-4">
      <div className="row justify-content-center">
        <div className="col-md-8">
          <div className="card">
            <div className="card-header">
              <h2 className="card-title mb-0">Новая транзакция</h2>
            </div>
            
            {error && (
              <div className="alert alert-danger m-3" role="alert">
                <i className="bi bi-exclamation-triangle-fill me-2"></i>
                {error}
              </div>
            )}

            <div className="card-body">
              <form onSubmit={handleSubmit}>
                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Тип пользователя</label>
                    <select
                      value={formData.user_type}
                      onChange={(e) => setFormData({...formData, user_type: e.target.value as 'individual' | 'legal'})}
                      className="form-select"
                      required
                    >
                      <option value="individual">Физическое лицо</option>
                      <option value="legal">Юридическое лицо</option>
                    </select>
                  </div>

                  <div className="col-md-6">
                    <label className="form-label">Тип транзакции</label>
                    <select
                      value={formData.trans_type}
                      onChange={(e) => setFormData({...formData, trans_type: e.target.value as 'credit' | 'debit'})}
                      className="form-select"
                      required
                    >
                      <option value="debit">Расход</option>
                      <option value="credit">Доход</option>
                    </select>
                  </div>
                </div>

                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Дата</label>
                    <input
                      type="date"
                      value={formData.date}
                      onChange={(e) => setFormData({...formData, date: e.target.value})}
                      className="form-control"
                      required
                    />
                  </div>

                  <div className="col-md-6">
                    <label className="form-label">Сумма</label>
                    <div className="input-group">
                      <input
                        type="number"
                        step="0.01"
                        value={formData.amount}
                        onChange={(e) => setFormData({...formData, amount: parseFloat(e.target.value)})}
                        className="form-control"
                        required
                        min="0"
                      />
                      <span className="input-group-text">₽</span>
                    </div>
                  </div>
                </div>

                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Категория</label>
                    <select
                      value={formData.category_id}
                      onChange={(e) => setFormData({...formData, category_id: parseInt(e.target.value)})}
                      className="form-select"
                      required
                    >
                      {categories.map(category => (
                        <option key={category.id} value={category.id}>
                          {category.name}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div className="col-md-6">
                    <label className="form-label">Статус</label>
                    <select
                      value={formData.status_id}
                      onChange={(e) => setFormData({...formData, status_id: parseInt(e.target.value)})}
                      className="form-select"
                      required
                    >
                      {statuses.map(status => (
                        <option key={status.id} value={status.id}>
                          {status.name}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Банк отправителя</label>
                    <select
                      value={formData.sender_bank}
                      onChange={(e) => setFormData({...formData, sender_bank: e.target.value})}
                      className="form-select"
                      required
                    >
                      <option value="">Выберите банк</option>
                      {banks.map(bank => (
                        <option key={bank.id} value={bank.name}>
                          {bank.name}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div className="col-md-6">
                    <label className="form-label">ИНН получателя</label>
                    <input
                      type="text"
                      value={formData.receiver_inn}
                      onChange={(e) => setFormData({...formData, receiver_inn: e.target.value})}
                      className="form-control"
                      required
                      pattern="[0-9]{10}|[0-9]{12}"
                      title="ИНН должен содержать 10 цифр для юр. лиц или 12 цифр для физ. лиц"
                    />
                  </div>
                </div>

                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Телефон получателя</label>
                    <input
                      type="tel"
                      value={formData.receiver_phone}
                      onChange={(e) => setFormData({...formData, receiver_phone: e.target.value})}
                      className="form-control"
                      required
                      pattern="[0-9]{11}"
                      title="Телефон должен содержать 11 цифр"
                    />
                  </div>
                </div>

                <div className="mb-3">
                  <label className="form-label">Комментарий</label>
                  <textarea
                    value={formData.comment}
                    onChange={(e) => setFormData({...formData, comment: e.target.value})}
                    rows={3}
                    className="form-control"
                    placeholder="Введите комментарий к транзакции..."
                  />
                </div>

                <div className="d-flex justify-content-end gap-2">
                  <button
                    type="button"
                    onClick={() => navigate('/transactions')}
                    className="btn btn-outline-secondary"
                  >
                    Отмена
                  </button>
                  <button
                    type="submit"
                    disabled={loading}
                    className="btn btn-primary"
                  >
                    {loading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        Создание...
                      </>
                    ) : (
                      'Создать транзакцию'
                    )}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TransactionForm; 