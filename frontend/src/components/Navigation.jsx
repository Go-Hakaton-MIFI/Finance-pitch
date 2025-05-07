import React from 'react';
import { Navbar, Nav, Container, Button } from 'react-bootstrap';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { logout } from '../store/authSlice';

const Navigation = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();
  const { user } = useSelector(state => state.auth);

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  const isActive = (path) => {
    return location.pathname === path;
  };

  return (
    <Navbar bg="white" expand="lg" className="shadow-sm">
      <Container fluid>
        <Navbar.Brand as={Link} to="/" className="text-primary">
          <i className="bi bi-wallet2 me-2"></i>
          Finance Manager
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link 
              as={Link} 
              to="/" 
              className={isActive('/') ? 'active' : ''}
            >
              <i className="bi bi-house-door me-1"></i>
              Главная
            </Nav.Link>
            <Nav.Link 
              as={Link} 
              to="/transactions" 
              className={isActive('/transactions') ? 'active' : ''}
            >
              <i className="bi bi-credit-card me-1"></i>
              Транзакции
            </Nav.Link>
            <Nav.Link 
              as={Link} 
              to="/analytics" 
              className={isActive('/analytics') ? 'active' : ''}
            >
              <i className="bi bi-graph-up me-1"></i>
              Аналитика
            </Nav.Link>
            <Nav.Link 
              as={Link} 
              to="/reports" 
              className={isActive('/reports') ? 'active' : ''}
            >
              <i className="bi bi-file-earmark-text me-1"></i>
              Отчеты
            </Nav.Link>
          </Nav>
          <Nav>
            <Nav.Link 
              as={Link} 
              to="/settings" 
              className={isActive('/settings') ? 'active' : ''}
            >
              <i className="bi bi-gear me-1"></i>
              Настройки
            </Nav.Link>
            <Nav.Link 
              as={Link} 
              to="/profile" 
              className={isActive('/profile') ? 'active' : ''}
            >
              <i className="bi bi-person-circle me-1"></i>
              {user?.firstName || 'Профиль'}
            </Nav.Link>
            <Button 
              variant="outline-danger" 
              onClick={handleLogout}
              className="ms-2"
            >
              <i className="bi bi-box-arrow-right me-1"></i>
              Выход
            </Button>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Navigation; 