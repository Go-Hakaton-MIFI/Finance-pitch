import { Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Transactions from './pages/Transactions';
import Categories from './pages/Categories';
import Reports from './pages/Reports';
import Profile from './pages/Profile';
import PrivateRoute from './components/PrivateRoute';

const routes = [
  {
    path: '/',
    element: <PrivateRoute><Dashboard /></PrivateRoute>
  },
  {
    path: '/login',
    element: <Login />
  },
  {
    path: '/register',
    element: <Register />
  },
  {
    path: '/transactions',
    element: <PrivateRoute><Transactions /></PrivateRoute>
  },
  {
    path: '/categories',
    element: <PrivateRoute><Categories /></PrivateRoute>
  },
  {
    path: '/reports',
    element: <PrivateRoute><Reports /></PrivateRoute>
  },
  {
    path: '/profile',
    element: <PrivateRoute><Profile /></PrivateRoute>
  },
  {
    path: '*',
    element: <Navigate to="/" replace />
  }
];

export default routes; 