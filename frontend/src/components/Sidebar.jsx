import React from 'react';
import { Nav } from 'react-bootstrap';
import { Link, useLocation } from 'react-router-dom';

const Sidebar = () => {
  const location = useLocation();

  const menuItems = [
    { path: '/dashboard', icon: 'bi-house-door', label: 'Панель управления' },
    { path: '/transactions', icon: 'bi-credit-card', label: 'Транзакции' },
    { path: '/analytics', icon: 'bi-graph-up', label: 'Аналитика' },
    { path: '/reports', icon: 'bi-file-earmark-text', label: 'Отчеты' },
    { path: '/profile', icon: 'bi-person', label: 'Профиль' },
  ];

  const handleLogout = () => {
    localStorage.removeItem('token');
    window.location.href = '/login';
  };

  return (
    <div className="sidebar" style={{ 
      width: '250px', 
      minHeight: '100vh',
      backgroundColor: '#f8f9fa',
      borderRight: '1px solid #e9ecef'
    }}>
      <div className="p-3">
        <h4 className="text-primary mb-4">Finance App</h4>
        <Nav className="flex-column">
          {menuItems.map((item) => (
            <Nav.Link
              key={item.path}
              as={Link}
              to={item.path}
              className={`d-flex align-items-center mb-2 ${
                location.pathname === item.path 
                  ? 'active bg-primary text-white' 
                  : 'text-primary hover-bg-light'
              }`}
              style={{
                borderRadius: '4px',
                transition: 'all 0.2s ease-in-out'
              }}
            >
              <i className={`bi ${item.icon} me-2`}></i>
              {item.label}
            </Nav.Link>
          ))}
          <Nav.Link
            onClick={handleLogout}
            className="d-flex align-items-center text-primary mt-4 hover-bg-light"
            style={{
              borderRadius: '4px',
              transition: 'all 0.2s ease-in-out'
            }}
          >
            <i className="bi bi-box-arrow-right me-2"></i>
            Выйти
          </Nav.Link>
        </Nav>
      </div>
    </div>
  );
};

export default Sidebar; 